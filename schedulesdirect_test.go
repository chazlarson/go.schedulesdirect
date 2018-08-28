package schedulesdirect

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// inspired by https://willnorris.com/2013/08/testing-in-go-github
// adapted from github.com/bbigras/go-schedulesdirect
var (
	mux    *http.ServeMux
	server *httptest.Server
	client *Client
)

func setup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	// schedules direct client configured to use test server
	client = &Client{
		BaseURL: fmt.Sprint(server.URL, "/"),
		HTTP:    http.DefaultClient,
		Token:   "d97c908ed44c25fdca302612c70584c8d5acd47a", // token1
	}
}

func testMethod(t *testing.T, r *http.Request, expectedMethod string) {
	t.Helper()
	if r.Method != expectedMethod {
		t.Fatalf("method (%s) != expectedMethod (%s)", r.Method, expectedMethod)
	}
}

func testHeader(t *testing.T, r *http.Request, header, expectedValue string) {
	t.Helper()
	if r.Header.Get(header) != expectedValue {
		t.Fatalf("token (%s) != expectedValue (%s)", r.Header.Get("token"), expectedValue)
	}
}

func testPayload(t *testing.T, r *http.Request, expect []byte) {
	t.Helper()
	data, errRead := ioutil.ReadAll(r.Body)
	if errRead != nil {
		t.Fatal(errRead)
	}

	if !bytes.Equal(data, expect) {
		t.Fatalf("payload doesn't match\nhas: >%s<\nexpect: >%s<", data, expect)
	}
}

func testUrlParameter(t *testing.T, r *http.Request, parameter, expectedValue string) {
	t.Helper()
	p := r.URL.Query().Get(parameter)

	if p != expectedValue {
		t.Fatalf("parameter (%s (%s)) != expectedValue (%s)", parameter, p, expectedValue)
	}
}

func ensureError(t *testing.T, err error, expectedCode ErrorCode) {
	t.Helper()
	if err == nil {
		t.Fatalf("was expecting error, did not get one")
	}

	if e, ok := err.(*BaseResponse); ok {
		if e.Code != expectedCode {
			t.Fatalf(`was expecting error to be of type "%s" (%d), was instead "%s" (%d)`, expectedCode.InternalCode(), int64(expectedCode), e.Code.InternalCode(), int64(expectedCode))
		}
		return
	} else {
		t.Fatalf("error was not of type BaseResponse, error string is: %s", err)
	}
}

func TestGetTokenOK(t *testing.T) {
	setup()

	mux.HandleFunc(fmt.Sprint("/", APIVersion, "/token"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "POST")

			var tokenResp TokenResponse

			errDecode := json.NewDecoder(r.Body).Decode(&tokenResp)
			if errDecode != nil {
				t.Fatal(errDecode)
			}

			fmt.Fprint(w, `{"code":0,"message":"OK","serverID":"serverID1","token":"d97c908ed44c25fdca302612c70584c8d5acd47a"}`)
		},
	)

	token, errToken := client.GetToken("user1", "pass1")
	if errToken != nil {
		t.Fatal(errToken)
	}

	if token != "d97c908ed44c25fdca302612c70584c8d5acd47a" {
		t.Fatalf("token doesn't match")
	}
}

func TestEncryptPassword(t *testing.T) {
	setup()

	hash, hashErr := encryptPassword("testpassword")
	if hashErr != nil {
		t.Fatal(hashErr)
	}
	if hash != "8bb6118f8fd6935ad0876a3be34a717d32708ffd" {
		t.Fail()
	}

}

func TestGetTokenInvalidUser(t *testing.T) {
	setup()

	mux.HandleFunc(fmt.Sprint("/", APIVersion, "/token"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "POST")

			var tokenResp TokenResponse

			errDecode := json.NewDecoder(r.Body).Decode(&tokenResp)
			if errDecode != nil {
				t.Fatal(errDecode)
			}

			fmt.Fprint(w, `{"response":"INVALID_USER","code":4003,"serverID":"serverID1","message":"Invalid user.","datetime":"2014-07-29T01:00:28Z"}`)
		},
	)

	_, errToken := client.GetToken("user1", "pass1")
	ensureError(t, errToken, ErrInvalidUser)
}

func TestGetStatusOK(t *testing.T) {
	setup()

	mux.HandleFunc(fmt.Sprint("/", APIVersion, "/status"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			testHeader(t, r, "token", "d97c908ed44c25fdca302612c70584c8d5acd47a")

			fmt.Fprint(w, `{"account":{"expires":"2014-09-26T19:07:28Z","messages":[],"maxLineups":4,"nextSuggestedConnectTime":"2014-07-29T22:43:22Z"},"lineups":[],"lastDataUpdate":"2014-07-28T14:48:59Z","notifications":[],"systemStatus":[{"date":"2012-12-17T16:24:47Z","status":"Online","details":"All servers running normally."}],"serverID":"serverID1","code":0}`)
		},
	)

	status, err := client.GetStatus()
	if err != nil {
		t.Fatal(err)
	}

	if len(status.SystemStatus) != 1 {
		t.Fail()
	} else if status.SystemStatus[0].Details != "All servers running normally." {
		t.Fail()
	}
}

func TestGetHeadendsOK(t *testing.T) {
	setup()

	mux.HandleFunc(fmt.Sprint("/", APIVersion, "/headends"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			testHeader(t, r, "token", "d97c908ed44c25fdca302612c70584c8d5acd47a")
			testUrlParameter(t, r, "country", "CAN")
			testUrlParameter(t, r, "postalcode", "H0H 0H0")

			fmt.Fprint(w, `[{"headend":"0000001","lineups":[{"name":"name1","uri":"uri1"},{"name":"name2","uri":"uri2"}],"location":"City1","transport":"type1"},{"headend":"0000002","lineups":[{"name":"name3","uri":"uri3"}],"location":"City2","transport":"type2"}]`)
		},
	)

	headends, errGetHeadends := client.GetHeadends("CAN", "H0H 0H0")
	if errGetHeadends != nil {
		t.Fatal(errGetHeadends)
	}

	if len(headends) != 2 {
		t.Fatalf("len(headends) != 2: %d", len(headends))
	} else {
		for _, headend := range headends {
			if headend.Headend == "0000001" {
				if len(headend.Lineups) != 2 {
					t.Fatalf(`len(headend.Lineups) != 2: %d`, len(headend.Lineups))
				} else if headend.Lineups[0].Name != "name1" {
					t.Fatalf(`headend.Lineups[0].Name != "name1": %s`, headend.Lineups[0].Name)
				}
			} else if headend.Headend == "0000002" {
				if len(headend.Lineups) != 1 {
					t.Fatalf(`len(headend.Lineups) != 1: %d`, len(headend.Lineups))
				}
			}
		}
	}
}

func TestGetHeadendsFailsWithMessage(t *testing.T) {
	setup()

	mux.HandleFunc(fmt.Sprint("/", APIVersion, "/headends"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			testHeader(t, r, "token", "d97c908ed44c25fdca302612c70584c8d5acd47a")
			testUrlParameter(t, r, "country", "CAN")
			testUrlParameter(t, r, "postalcode", "H0H 0H0")

			fmt.Fprint(w, `{"response":"INVALID_PARAMETER:COUNTRY","code":2050,"serverID":"serverID1","message":"The COUNTRY parameter must be ISO-3166-1 alpha 3. See http:\/\/en.wikipedia.org\/wiki\/ISO_3166-1_alpha-3","datetime":"2014-07-29T23:16:52Z"}`)
		},
	)

	_, errGetHeadends := client.GetHeadends("CAN", "H0H 0H0")
	ensureError(t, errGetHeadends, ErrInvalidParameterCountry)
}

func TestGetHeadendsFailsWithMessage2(t *testing.T) {
	setup()

	mux.HandleFunc(fmt.Sprint("/", APIVersion, "/headends"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			testHeader(t, r, "token", "d97c908ed44c25fdca302612c70584c8d5acd47a")
			testUrlParameter(t, r, "country", "CAN")
			testUrlParameter(t, r, "postalcode", "H0H 0H0")

			fmt.Fprint(w, `{"response":"REQUIRED_PARAMETER_MISSING:COUNTRY","code":2004,"serverID":"serverID1","message":"In order to search for lineups, you must supply a 3-letter country parameter.","datetime":"2014-07-29T23:15:18Z"}`)
		},
	)

	_, errGetHeadends := client.GetHeadends("CAN", "H0H 0H0")
	ensureError(t, errGetHeadends, ErrRequiredParameterMissingCountry)
}

func TestAddLineupOK(t *testing.T) {
	setup()

	mux.HandleFunc(fmt.Sprint("/", APIVersion, "/lineups/CAN-0000001-X"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "PUT")
			testHeader(t, r, "token", "d97c908ed44c25fdca302612c70584c8d5acd47a")
			fmt.Fprint(w, `{"response":"OK","code":0,"serverID":"serverID1","message":"Added lineup.","changesRemaining":5,"datetime":"2014-07-30T01:50:59Z"}`)
		},
	)

	changeResp, errAddLineup := client.AddLineup("CAN-0000001-X")
	if errAddLineup != nil {
		t.Fatal(errAddLineup)
	}

	if changeResp.ChangesRemaining != 5 {
		t.Fail()
	}
}

func TestAddLineupFailsDuplicate(t *testing.T) {
	setup()

	mux.HandleFunc(fmt.Sprint("/", APIVersion, "/lineups/CAN-0000001-X"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "PUT")
			testHeader(t, r, "token", "d97c908ed44c25fdca302612c70584c8d5acd47a")
			fmt.Fprint(w, `{"response":"DUPLICATE_HEADEND","code":2100,"serverID":"serverID1","message":"Headend already in account.","datetime":"2014-07-30T02:01:37Z"}`)
		},
	)

	_, errAddLineup := client.AddLineup("CAN-0000001-X")
	ensureError(t, errAddLineup, ErrDuplicateLineup)
}

func TestAddLineupFailsInvalidLineup(t *testing.T) {
	setup()

	mux.HandleFunc(fmt.Sprint("/", APIVersion, "/lineups/CAN-0000001-X"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "PUT")
			testHeader(t, r, "token", "d97c908ed44c25fdca302612c70584c8d5acd47a")

			fmt.Fprint(w, `{"response":"INVALID_LINEUP","code":2105,"serverID":"serverID1","message":"The lineup you submitted doesn't exist.","datetime":"2014-07-30T02:02:04Z"}`)
		},
	)

	_, errAddLineup := client.AddLineup("CAN-0000001-X")
	ensureError(t, errAddLineup, ErrInvalidLineup)
}

func TestAddLineupFailsInvalidUser(t *testing.T) {
	setup()

	mux.HandleFunc(fmt.Sprint("/", APIVersion, "/lineups/CAN-0000001-X"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "PUT")
			testHeader(t, r, "token", "d97c908ed44c25fdca302612c70584c8d5acd47a")
			fmt.Fprint(w, `{"response":"INVALID_USER","code":4003,"serverID":"serverID1","message":"Invalid user.","datetime":"2014-07-30T01:48:11Z"}`)
		},
	)

	_, errAddLineup := client.AddLineup("CAN-0000001-X")
	ensureError(t, errAddLineup, ErrInvalidUser)
}

func TestDeleteLineupOK(t *testing.T) {
	setup()

	mux.HandleFunc(fmt.Sprint("/", APIVersion, "/lineups/CAN-0000001-X"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "DELETE")
			testHeader(t, r, "token", "d97c908ed44c25fdca302612c70584c8d5acd47a")
			fmt.Fprint(w, `{"response":"OK","code":0,"serverID":"serverid1","message":"Deleted lineup.","changesRemaining":"5","datetime":"2014-07-30T03:27:23Z"}`)
		},
	)

	changeResp, errDeleteLineup := client.DeleteLineup("CAN-0000001-X")
	if errDeleteLineup != nil {
		t.Fatal(errDeleteLineup)
	}

	if changeResp.ChangesRemaining != 5 {
		t.Fail()
	}
}

func TestDeleteLineupFailsInvalidLineup(t *testing.T) {
	setup()

	mux.HandleFunc(fmt.Sprint("/", APIVersion, "/lineups/CAN-0000001-X"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "DELETE")
			testHeader(t, r, "token", "d97c908ed44c25fdca302612c70584c8d5acd47a")

			fmt.Fprint(w, `{"response":"INVALID_LINEUP","code":2105,"serverID":"serverID1","message":"The lineup you submitted doesn't exist.","datetime":"2014-07-30T02:02:04Z"}`)
		},
	)

	_, errDeleteLineup := client.DeleteLineup("CAN-0000001-X")
	ensureError(t, errDeleteLineup, ErrInvalidLineup)
}

func TestGetLineupsOK(t *testing.T) {
	setup()

	mux.HandleFunc(fmt.Sprint("/", APIVersion, "/lineups"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			testHeader(t, r, "token", "d97c908ed44c25fdca302612c70584c8d5acd47a")

			fmt.Fprint(w, `{"serverID":"serverid1","datetime":"2014-07-30T02:34:37Z","lineups":[{"name":"name1","type":"type1","location":"location1","uri":"uri1"}]}`)
		},
	)

	lineups, errGetLineups := client.GetLineups()
	if errGetLineups != nil {
		t.Fatal(errGetLineups)
	}

	if len(lineups.Lineups) != 1 {
		t.Fatalf("len(lineups.Lineups) != 1: %d", len(lineups.Lineups))
	} else if lineups.Lineups[0].Name != "name1" {
		t.Fatalf(`lineups.Lineups[0].Name != "name1": %s`, lineups.Lineups[0].Name)
	}
}

func TestGetLineupsFailsNoHeadends(t *testing.T) {
	setup()

	mux.HandleFunc(fmt.Sprint("/", APIVersion, "/lineups"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			testHeader(t, r, "token", "d97c908ed44c25fdca302612c70584c8d5acd47a")

			// bug with the web service?
			http.Error(w, "", http.StatusBadRequest)

			fmt.Fprint(w, `{"response":"NO_LINEUPS","code":4102,"serverID":"serverID1","message":"No lineups have been added to this account.","datetime":"2014-07-30T01:21:56Z"}`)
		},
	)

	_, errGetLineups := client.GetLineups()
	ensureError(t, errGetLineups, ErrNoLineups)
}

func TestScheduleEquality(t *testing.T) {
	expected := `[{"stationID":"20454","programs":[{"programID":"SH005371070000","airDateTime":"2015-03-03T00:00:00Z","duration":1800,"md5":"Sy8HEMBPcuiAx3FBukUhKQ","new":true,"audioProperties":["stereo","cc"],"videoProperties":["hdtv"]},{"programID":"EP000014577244","airDateTime":"2015-03-03T00:30:00Z","duration":1800,"md5":"25DNXVXO192JI7Y9vSW9lQ","new":true,"audioProperties":["stereo","cc"],"videoProperties":["hdtv"]},{"programID":"EP014145320829","airDateTime":"2015-03-03T23:30:00Z","duration":1800,"md5":"A6RCPnx4SjKN3oaZtXxNfw","new":true,"audioProperties":["stereo","cc"],"videoProperties":["hdtv"]}],"metadata":{"modified":"2015-03-02T15:56:02Z","md5":"UtL+hq0sqtCTZVFrGHZ5sg","startDate":"2015-03-03"}},{"stationID":"20454","programs":[{"programID":"SH005371070000","airDateTime":"2015-03-04T00:00:00Z","duration":1800,"md5":"Sy8HEMBPcuiAx3FBukUhKQ","new":true,"audioProperties":["stereo","cc"],"videoProperties":["hdtv"]},{"programID":"EP014145320830","airDateTime":"2015-03-04T23:30:00Z","duration":1800,"md5":"LdPQVqsQJWTcX1b5k2VPGQ","new":true,"audioProperties":["stereo","cc"],"videoProperties":["hdtv"]}],"metadata":{"modified":"2015-03-02T15:56:02Z","md5":"ZZPts55w9WUP1rMRvKsGDw","startDate":"2015-03-04"}},{"stationID":"10021","programs":[{"programID":"EP018632100004","airDateTime":"2015-03-03T01:56:00Z","duration":3840,"md5":"J+AOJ/ofAQdp12Bh3U+C+A","audioProperties":["cc"],"ratings":[{"body":"USA Parental Rating","code":"TV14"}]}]}]`

	schedule := make([]Schedule, 0)

	if unmarshalErr := json.Unmarshal([]byte(expected), &schedule); unmarshalErr != nil {
		t.Fatalf("error when unmarshalling expected json to slice of schedule: %s", unmarshalErr)
	}

	marshalled, marshalErr := json.Marshal(schedule)
	if marshalErr != nil {
		t.Fatalf("error marshalling slice of schedule to json: %s", marshalErr)
	}

	assert.JSONEq(t, expected, string(marshalled))
}

func TestGetChannelsOK(t *testing.T) {
	setup()

	mux.HandleFunc(fmt.Sprint("/", APIVersion, "/lineups/CAN-0000001-X"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			testHeader(t, r, "token", "d97c908ed44c25fdca302612c70584c8d5acd47a")
			fmt.Fprint(w, `{"map": [{"channel": "101","stationID": "10001"},{"channel": "1933","stationID": "10001"}],"metadata": {"lineup": "CAN-0000000-X","modified": "2014-07-29T16:38:09Z","transport": "transport1"},"stations": [{"affiliate": "affiliate1","broadcaster": {"city": "Unknown","country": "Unknown","postalcode": "00000"},"callsign": "callsign1","language": "en","name": "name1","stationID": "10001"},       {"callsign": "callsign2","language": "en","logo": {"URL": "https://domain/path/file.png","dimension": "w=360px|h=270px","md5": "ba5b5b5085baac6da247564039c03c9e"},"name": "name2","stationID": "10002"}]}`)
		},
	)

	channelMapping, errGetChannels := client.GetChannels("CAN-0000001-X", false)
	if errGetChannels != nil {
		t.Fatal(errGetChannels)
	}

	if len(channelMapping.Map) != 2 {
		t.Fail()
	}
	if len(channelMapping.Stations) != 2 {
		t.Fail()
	}
	if channelMapping.Metadata.Lineup != "CAN-0000000-X" {
		t.Fail()
	}
}

func TestGetChannelsFailsLineupNotFound(t *testing.T) {
	setup()

	mux.HandleFunc(fmt.Sprint("/", APIVersion, "/lineups/CAN-0000001-X"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			testHeader(t, r, "token", "d97c908ed44c25fdca302612c70584c8d5acd47a")
			fmt.Fprint(w, `{"response":"LINEUP_NOT_FOUND","code":2101,"serverID":"serverid1","message":"Lineup not in account. Add lineup to account before requesting mapping.","datetime":"2014-07-30T04:14:27Z"}`)
		},
	)

	_, errGetChannels := client.GetChannels("CAN-0000001-X", false)
	ensureError(t, errGetChannels, ErrLineupNotFound)
}

func TestProgramInfoEquality(t *testing.T) {
	expected := `[{"programID":"EP000000060003","titles":[{"title120":"'Allo 'Allo!"}],"eventDetails":{"subType":"Series"},"descriptions":{"description1000":[{"descriptionLanguage":"en","description":"A disguised British Intelligence officer is sent to help the airmen."}]},"originalAirDate":"1985-11-04","genres":["Sitcom"],"episodeTitle150":"The Poloceman Cometh","metadata":[{"Gracenote":{"season":2,"episode":3}}],"cast":[{"personId":"383774","nameId":"392649","name":"Gorden Kaye","role":"Actor","billingOrder":"01"},{"personId":"246840","nameId":"250387","name":"Carmen Silvera","role":"Actor","billingOrder":"02"},{"personId":"376955","nameId":"385830","name":"Rose Hill","role":"Actor","billingOrder":"03"},{"personId":"259773","nameId":"263340","name":"Vicki Michelle","role":"Actor","billingOrder":"04"},{"personId":"353113","nameId":"361987","name":"Kirsten Cooke","role":"Actor","billingOrder":"05"},{"personId":"77787","nameId":"77787","name":"Richard Marner","role":"Actor","billingOrder":"06"},{"personId":"230921","nameId":"234193","name":"Guy Siner","role":"Actor","billingOrder":"07"},{"personId":"374934","nameId":"383809","name":"Kim Hartman","role":"Actor","billingOrder":"08"},{"personId":"369151","nameId":"378026","name":"Richard Gibson","role":"Actor","billingOrder":"09"},{"personId":"343690","nameId":"352564","name":"Arthur Bostrom","role":"Actor","billingOrder":"10"},{"personId":"352557","nameId":"361431","name":"John D. Collins","role":"Actor","billingOrder":"11"},{"personId":"605275","nameId":"627734","name":"Nicholas Frankau","role":"Actor","billingOrder":"12"},{"personId":"373394","nameId":"382269","name":"Jack Haig","role":"Actor","billingOrder":"13"}],"crew":[{"personId":"354407","nameId":"363281","name":"David Croft","role":"Director","billingOrder":"01"},{"personId":"354407","nameId":"363281","name":"David Croft","role":"Writer","billingOrder":"02"},{"personId":"105145","nameId":"105145","name":"Jeremy Lloyd","role":"Writer","billingOrder":"03"}],"showType":"Series","hasImageArtwork":true,"md5":"Jo5NKxoo44xRvBCAq8QT2A"},{"programID":"EP000000510142","titles":[{"title120":"A Different World"}],"eventDetails":{"subType":"Series"},"descriptions":{"description1000":[{"descriptionLanguage":"en","description":"Whitley and Dwayne tell new students about their honeymoon in Los Angeles."}]},"originalAirDate":"1992-09-24","genres":["Sitcom"],"episodeTitle150":"Honeymoon in L.A.","metadata":[{"Gracenote":{"season":6,"episode":1}}],"cast":[{"personId":"700","nameId":"700","name":"Jasmine Guy","role":"Actor","billingOrder":"01"},{"personId":"729","nameId":"729","name":"Kadeem Hardison","role":"Actor","billingOrder":"02"},{"personId":"120","nameId":"120","name":"Darryl M. Bell","role":"Actor","billingOrder":"03"},{"personId":"1729","nameId":"1729","name":"Cree Summer","role":"Actor","billingOrder":"04"},{"personId":"217","nameId":"217","name":"Charnele Brown","role":"Actor","billingOrder":"05"},{"personId":"1811","nameId":"1811","name":"Glynn Turman","role":"Actor","billingOrder":"06"},{"personId":"1232","nameId":"1232","name":"Lou Myers","role":"Actor","billingOrder":"07"},{"personId":"1363","nameId":"1363","name":"Jada Pinkett","role":"Guest Star","billingOrder":"08"},{"personId":"222967","nameId":"225536","name":"Ajai Sanders","role":"Guest Star","billingOrder":"09"},{"personId":"181744","nameId":"183292","name":"Karen Malina White","role":"Guest Star","billingOrder":"10"},{"personId":"305017","nameId":"318897","name":"Patrick Y. Malone","role":"Guest Star","billingOrder":"11"},{"personId":"9841","nameId":"9841","name":"Bumper Robinson","role":"Guest Star","billingOrder":"12"},{"personId":"426422","nameId":"435297","name":"Sister Souljah","role":"Guest Star","billingOrder":"13"},{"personId":"25","nameId":"25","name":"Debbie Allen","role":"Guest Star","billingOrder":"14"},{"personId":"668","nameId":"668","name":"Gilbert Gottfried","role":"Guest Star","billingOrder":"15"}],"showType":"Series","hasImageArtwork":true,"md5":"P5kz0QmCeYxIA+yL0H4DWw"}]`

	programs := make([]ProgramInfo, 0)

	if unmarshalErr := json.Unmarshal([]byte(expected), &programs); unmarshalErr != nil {
		t.Fatalf("error when unmarshalling expected json to slice of programinfo: %s", unmarshalErr)
	}

	marshalled, marshalErr := json.Marshal(programs)
	if marshalErr != nil {
		t.Fatalf("error marshalling slice of programs to json: %s", marshalErr)
	}

	assert.JSONEq(t, expected, string(marshalled))
}

func TestGetProgramInfoOK(t *testing.T) {
	setup()

	mux.HandleFunc(fmt.Sprint("/", APIVersion, "/programs"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "POST")
			testHeader(t, r, "token", "d97c908ed44c25fdca302612c70584c8d5acd47a")
			testPayload(t, r, []byte(`["program1","program2"]`))

			fmt.Fprint(w, `[{"programID":"program1","titles":[{"title120":"title1"}],"eventDetails":{"subType":"subType1"},"originalAirDate":"2012-01-01","genres":["genre1"],"showType":"type1","md5":"edbb1c792032ba8685fd021c28c6ea74"},
{"programID":"program2","titles":[{"title120":"title2"}],"eventDetails":{"subType":"subType2"},"originalAirDate":"2012-01-01","genres":["genre2"],"showType":"type2","md5":"edbb1c792032ba8685fd021c28c6ea74"}]`)
		},
	)

	programs, err := client.GetProgramInfo([]string{
		"program1",
		"program2",
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(programs) != 2 {
		t.Fail()
	} else {
		if programs[0].ProgramID != "program1" {
			t.Fail()
		}
		if programs[1].ProgramID != "program2" {
			t.Fail()
		}
	}
}

func TestGetProgramInfoFailsRequiredRequestMissing(t *testing.T) {
	setup()

	mux.HandleFunc(fmt.Sprint("/", APIVersion, "/programs"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "POST")
			testHeader(t, r, "token", "d97c908ed44c25fdca302612c70584c8d5acd47a")
			testPayload(t, r, []byte(`["program1","program2"]`))

			fmt.Fprint(w, `{"response":"REQUIRED_REQUEST_MISSING","code":2002,"serverID":"serverid1","message":"Did not receive request.","datetime":"2014-07-30T05:02:22Z"}`)
		},
	)

	_, err := client.GetProgramInfo([]string{
		"program1",
		"program2",
	})
	ensureError(t, err, ErrRequiredRequestMissing)
}

func TestGetProgramInfoFailsDeflateRequired(t *testing.T) {
	setup()

	mux.HandleFunc(fmt.Sprint("/", APIVersion, "/programs"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "POST")
			testHeader(t, r, "token", "d97c908ed44c25fdca302612c70584c8d5acd47a")
			testPayload(t, r, []byte(`["program1","program2"]`))

			fmt.Fprint(w, `{"response":"DEFLATE_REQUIRED","code":1002,"serverID":"serverid1","message":"Did not receive Accept-Encoding: deflate in request","datetime":"2014-07-30T05:02:42Z"}`)
		},
	)

	_, err := client.GetProgramInfo([]string{
		"program1",
		"program2",
	})
	ensureError(t, err, ErrDeflateRequired)
}

func TestGetSchedulesOK(t *testing.T) {
	setup()

	mux.HandleFunc(fmt.Sprint("/", APIVersion, "/schedules"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "POST")
			testHeader(t, r, "token", "d97c908ed44c25fdca302612c70584c8d5acd47a")
			testPayload(t, r, []byte(`[{"stationID":"10001"},{"stationID":"10002"}]`))

			fmt.Fprint(w, `[{"metadata": {"endDate": "2014-08-12","startDate": "2014-07-30"},"programs": [{"airDateTime": "2014-07-30T00:30:00Z","audioProperties": ["ap1","ap2"],"contentRating": [{"body": "body1","code": "code1"}],"duration": 1800,"md5": "exubfjxJmKcSe52dVLj83g","new": true,"programID": "program1","syndication": {"source": "ss1","type": "st1"}},{"airDateTime": "2014-08-12T23:30:00Z","audioProperties": ["ap3","ap4","ap5"],"contentAdvisory": {"rating1": ["stuff1","stuff2"]},"contentRating": [{"body": "body2","code": "code2"}],"duration": 1800,"md5": "5BxxvnI4Nv9ZuT9oQvOpQA","programID": "program2","syndication": {"source": "ss2","type": "st2"}}],"stationID": "10001"},
{"metadata": {"endDate": "2014-08-12","startDate": "2014-07-30"},"programs": [{"airDateTime": "2014-07-30T00:30:00Z","duration": 1800,"md5": "exubfjxJmKcSe52dVLj83g","new": true,"programID": "program3","syndication": {"source": "ss3","type": "st3"}}],"stationID": "10002"}]`)
		},
	)

	schedules, err := client.GetSchedules([]StationScheduleRequest{
		StationScheduleRequest{StationID: "10001"},
		StationScheduleRequest{StationID: "10002"},
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(schedules) != 2 {
		t.Fail()
	} else {
		if schedules[0].StationID != "10001" {
			t.Fail()
		}

		if len(schedules[0].Programs) != 2 {
			t.Fail()
		}
		if len(schedules[1].Programs) != 1 {
			t.Fail()
		}

		if schedules[1].StationID != "10002" {
			t.Fail()
		}
	}
}

func TestGetSchedulesFailsStationNotInLineup(t *testing.T) {
	setup()

	mux.HandleFunc(fmt.Sprint("/", APIVersion, "/schedules"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "POST")
			testHeader(t, r, "token", "d97c908ed44c25fdca302612c70584c8d5acd47a")
			testPayload(t, r, []byte(`[{"stationID":"10002"}]`))

			fmt.Fprint(w, `{"stationID":10002,"response":"ERROR","code":2200,"serverID":"serverid1","message":"This stationID (10002) is not in any of your lineups.","datetime":"2014-07-30T17:14:56Z"}`)
		},
	)

	_, err := client.GetSchedules([]StationScheduleRequest{
		StationScheduleRequest{StationID: "10002"},
	})
	ensureError(t, err, ErrStationIDNotFound)
}
