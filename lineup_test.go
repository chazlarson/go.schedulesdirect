package schedulesdirect

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestGetHeadendsOK(t *testing.T) {
	setup()

	mux.HandleFunc(fmt.Sprint("/", APIVersion, "/headends"),
		func(w http.ResponseWriter, r *http.Request) {
			ensureMethod(t, r, "GET")
			ensureHeader(t, r, "token", "d97c908ed44c25fdca302612c70584c8d5acd47a")
			ensureURLParameter(t, r, "country", "CAN")
			ensureURLParameter(t, r, "postalcode", "H0H 0H0")

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
			ensureMethod(t, r, "GET")
			ensureHeader(t, r, "token", "d97c908ed44c25fdca302612c70584c8d5acd47a")
			ensureURLParameter(t, r, "country", "CAN")
			ensureURLParameter(t, r, "postalcode", "H0H 0H0")

			fmt.Fprint(w, getBaseResponse(ErrInvalidParameterCountry))
		},
	)

	_, errGetHeadends := client.GetHeadends("CAN", "H0H 0H0")
	ensureError(t, errGetHeadends, ErrInvalidParameterCountry)
}

func TestGetHeadendsFailsWithMessage2(t *testing.T) {
	setup()

	mux.HandleFunc(fmt.Sprint("/", APIVersion, "/headends"),
		func(w http.ResponseWriter, r *http.Request) {
			ensureMethod(t, r, "GET")
			ensureHeader(t, r, "token", "d97c908ed44c25fdca302612c70584c8d5acd47a")
			ensureURLParameter(t, r, "country", "CAN")
			ensureURLParameter(t, r, "postalcode", "H0H 0H0")

			fmt.Fprint(w, getBaseResponse(ErrRequiredParameterMissingCountry))
		},
	)

	_, errGetHeadends := client.GetHeadends("CAN", "H0H 0H0")
	ensureError(t, errGetHeadends, ErrRequiredParameterMissingCountry)
}

func TestAddLineupOK(t *testing.T) {
	setup()

	mux.HandleFunc(fmt.Sprint("/", APIVersion, "/lineups/CAN-0000001-X"),
		func(w http.ResponseWriter, r *http.Request) {
			ensureMethod(t, r, "PUT")
			ensureHeader(t, r, "token", "d97c908ed44c25fdca302612c70584c8d5acd47a")

			baseResp := getBaseResponse(ErrOK)

			fmt.Fprintf(w, fmt.Sprintf(`%s, "changesRemaining": "5"}`, baseResp[:len(baseResp)-1]))
		},
	)

	changeResp, errAddLineup := client.AddLineup("CAN-0000001-X")
	if errAddLineup != nil {
		t.Fatal(errAddLineup)
	}

	if changeResp.ChangesRemaining != 5 {
		t.FailNow()
	}
}

func TestAddLineupFailsDuplicate(t *testing.T) {
	setup()

	mux.HandleFunc(fmt.Sprint("/", APIVersion, "/lineups/CAN-0000001-X"),
		func(w http.ResponseWriter, r *http.Request) {
			ensureMethod(t, r, "PUT")
			ensureHeader(t, r, "token", "d97c908ed44c25fdca302612c70584c8d5acd47a")

			fmt.Fprint(w, getBaseResponse(ErrDuplicateLineup))
		},
	)

	_, errAddLineup := client.AddLineup("CAN-0000001-X")
	ensureError(t, errAddLineup, ErrDuplicateLineup)
}

func TestAddLineupFailsInvalidLineup(t *testing.T) {
	setup()

	mux.HandleFunc(fmt.Sprint("/", APIVersion, "/lineups/CAN-0000001-X"),
		func(w http.ResponseWriter, r *http.Request) {
			ensureMethod(t, r, "PUT")
			ensureHeader(t, r, "token", "d97c908ed44c25fdca302612c70584c8d5acd47a")

			fmt.Fprint(w, getBaseResponse(ErrInvalidLineup))
		},
	)

	_, errAddLineup := client.AddLineup("CAN-0000001-X")
	ensureError(t, errAddLineup, ErrInvalidLineup)
}

func TestAddLineupFailsInvalidUser(t *testing.T) {
	setup()

	mux.HandleFunc(fmt.Sprint("/", APIVersion, "/token"),
		func(w http.ResponseWriter, r *http.Request) {
			ensureMethod(t, r, "POST")

			var tokenResp TokenResponse

			errDecode := json.NewDecoder(r.Body).Decode(&tokenResp)
			if errDecode != nil {
				t.Fatal(errDecode)
			}

			baseResp := getBaseResponse(ErrOK)

			fmt.Fprintf(w, fmt.Sprintf(`%s, "token": "d97c908ed44c25fdca302612c70584c8d5acd47a"}`, baseResp[:len(baseResp)-1]))
		},
	)

	mux.HandleFunc(fmt.Sprint("/", APIVersion, "/lineups/CAN-0000001-X"),
		func(w http.ResponseWriter, r *http.Request) {
			ensureMethod(t, r, "PUT")
			ensureHeader(t, r, "token", "d97c908ed44c25fdca302612c70584c8d5acd47a")

			fmt.Fprint(w, getBaseResponse(ErrInvalidUser))
		},
	)

	_, errAddLineup := client.AddLineup("CAN-0000001-X")
	ensureError(t, errAddLineup, ErrInvalidUser)
}

func TestDeleteLineupOK(t *testing.T) {
	setup()

	mux.HandleFunc(fmt.Sprint("/", APIVersion, "/lineups/CAN-0000001-X"),
		func(w http.ResponseWriter, r *http.Request) {
			ensureMethod(t, r, "DELETE")
			ensureHeader(t, r, "token", "d97c908ed44c25fdca302612c70584c8d5acd47a")

			baseResp := getBaseResponse(ErrOK)

			fmt.Fprintf(w, fmt.Sprintf(`%s, "changesRemaining": "5"}`, baseResp[:len(baseResp)-1]))
		},
	)

	changeResp, errDeleteLineup := client.DeleteLineup("CAN-0000001-X")
	if errDeleteLineup != nil {
		t.Fatal(errDeleteLineup)
	}

	if changeResp.ChangesRemaining != 5 {
		t.FailNow()
	}
}

func TestDeleteLineupFailsInvalidLineup(t *testing.T) {
	setup()

	mux.HandleFunc(fmt.Sprint("/", APIVersion, "/lineups/CAN-0000001-X"),
		func(w http.ResponseWriter, r *http.Request) {
			ensureMethod(t, r, "DELETE")
			ensureHeader(t, r, "token", "d97c908ed44c25fdca302612c70584c8d5acd47a")

			fmt.Fprint(w, getBaseResponse(ErrInvalidLineup))
		},
	)

	_, errDeleteLineup := client.DeleteLineup("CAN-0000001-X")
	ensureError(t, errDeleteLineup, ErrInvalidLineup)
}

func TestGetLineupsOK(t *testing.T) {
	setup()

	mux.HandleFunc(fmt.Sprint("/", APIVersion, "/lineups"),
		func(w http.ResponseWriter, r *http.Request) {
			ensureMethod(t, r, "GET")
			ensureHeader(t, r, "token", "d97c908ed44c25fdca302612c70584c8d5acd47a")

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
			ensureMethod(t, r, "GET")
			ensureHeader(t, r, "token", "d97c908ed44c25fdca302612c70584c8d5acd47a")

			fmt.Fprint(w, getBaseResponse(ErrNoLineups))
		},
	)

	_, errGetLineups := client.GetLineups()
	ensureError(t, errGetLineups, ErrNoLineups)
}

func TestGetChannelsOK(t *testing.T) {
	setup()

	mux.HandleFunc(fmt.Sprint("/", APIVersion, "/lineups/CAN-0000001-X"),
		func(w http.ResponseWriter, r *http.Request) {
			ensureMethod(t, r, "GET")
			ensureHeader(t, r, "token", "d97c908ed44c25fdca302612c70584c8d5acd47a")
			fmt.Fprint(w, `{"map": [{"channel": "101","stationID": "10001"},{"channel": "1933","stationID": "10001"}],"metadata": {"lineup": "CAN-0000000-X","modified": "2014-07-29T16:38:09Z","transport": "transport1"},"stations": [{"affiliate": "affiliate1","broadcaster": {"city": "Unknown","country": "Unknown","postalcode": "00000"},"callsign": "callsign1","language": "en","name": "name1","stationID": "10001"},       {"callsign": "callsign2","language": "en","logo": {"URL": "https://domain/path/file.png","dimension": "w=360px|h=270px","md5": "ba5b5b5085baac6da247564039c03c9e"},"name": "name2","stationID": "10002"}]}`)
		},
	)

	channelMapping, errGetChannels := client.GetChannels("CAN-0000001-X", false)
	if errGetChannels != nil {
		t.Fatal(errGetChannels)
	}

	if len(channelMapping.Map) != 2 {
		t.FailNow()
	}
	if len(channelMapping.Stations) != 2 {
		t.FailNow()
	}
	if channelMapping.Metadata.Lineup != "CAN-0000000-X" {
		t.FailNow()
	}
}

func TestGetChannelsFailsLineupNotFound(t *testing.T) {
	setup()

	mux.HandleFunc(fmt.Sprint("/", APIVersion, "/lineups/CAN-0000001-X"),
		func(w http.ResponseWriter, r *http.Request) {
			ensureMethod(t, r, "GET")
			ensureHeader(t, r, "token", "d97c908ed44c25fdca302612c70584c8d5acd47a")

			fmt.Fprint(w, getBaseResponse(ErrLineupNotFound))
		},
	)

	_, errGetChannels := client.GetChannels("CAN-0000001-X", false)
	ensureError(t, errGetChannels, ErrLineupNotFound)
}
