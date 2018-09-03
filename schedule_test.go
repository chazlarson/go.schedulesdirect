package schedulesdirect

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestGetSchedulesOK(t *testing.T) {
	setup()

	mux.HandleFunc(fmt.Sprint("/", APIVersion, "/schedules"),
		func(w http.ResponseWriter, r *http.Request) {
			ensureMethod(t, r, "POST")
			ensureHeader(t, r, "token", "d97c908ed44c25fdca302612c70584c8d5acd47a")
			ensurePayload(t, r, []byte(`[{"stationID":"10001"},{"stationID":"10002"}]`))

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
		t.FailNow()
	} else {
		if schedules[0].StationID != "10001" {
			t.FailNow()
		}

		if len(schedules[0].Programs) != 2 {
			t.FailNow()
		}
		if len(schedules[1].Programs) != 1 {
			t.FailNow()
		}

		if schedules[1].StationID != "10002" {
			t.FailNow()
		}
	}
}

func TestGetSchedulesFailsStationNotInLineup(t *testing.T) {
	setup()

	mux.HandleFunc(fmt.Sprint("/", APIVersion, "/schedules"),
		func(w http.ResponseWriter, r *http.Request) {
			ensureMethod(t, r, "POST")
			ensureHeader(t, r, "token", "d97c908ed44c25fdca302612c70584c8d5acd47a")
			ensurePayload(t, r, []byte(`[{"stationID":"10002"}]`))

			fmt.Fprint(w, getBaseResponse(ErrStationIDNotFound))
		},
	)

	_, err := client.GetSchedules([]StationScheduleRequest{
		StationScheduleRequest{StationID: "10002"},
	})
	ensureError(t, err, ErrStationIDNotFound)
}
