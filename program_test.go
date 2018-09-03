package schedulesdirect

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			ensureMethod(t, r, "POST")
			ensureHeader(t, r, "token", "d97c908ed44c25fdca302612c70584c8d5acd47a")
			ensurePayload(t, r, []byte(`["program1","program2"]`))

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
		t.FailNow()
	} else {
		if programs[0].ProgramID != "program1" {
			t.FailNow()
		}
		if programs[1].ProgramID != "program2" {
			t.FailNow()
		}
	}
}

func TestGetProgramInfoFailsRequiredRequestMissing(t *testing.T) {
	setup()

	mux.HandleFunc(fmt.Sprint("/", APIVersion, "/programs"),
		func(w http.ResponseWriter, r *http.Request) {
			ensureMethod(t, r, "POST")
			ensureHeader(t, r, "token", "d97c908ed44c25fdca302612c70584c8d5acd47a")
			ensurePayload(t, r, []byte(`["program1","program2"]`))

			fmt.Fprint(w, getBaseResponse(ErrRequiredRequestMissing))
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
			ensureMethod(t, r, "POST")
			ensureHeader(t, r, "token", "d97c908ed44c25fdca302612c70584c8d5acd47a")
			ensurePayload(t, r, []byte(`["program1","program2"]`))

			fmt.Fprint(w, getBaseResponse(ErrDeflateRequired))
		},
	)

	_, err := client.GetProgramInfo([]string{
		"program1",
		"program2",
	})
	ensureError(t, err, ErrDeflateRequired)
}
