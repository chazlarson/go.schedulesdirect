package schedulesdirect

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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
		BaseURL:        fmt.Sprint(server.URL, "/"),
		HTTP:           http.DefaultClient,
		Token:          "d97c908ed44c25fdca302612c70584c8d5acd47a", // token1
		TokenExpiresAt: time.Now().Add(24 * time.Hour),
	}
}

func ensureMethod(t *testing.T, r *http.Request, expectedMethod string) {
	t.Helper()
	if r.Method != expectedMethod {
		t.Fatalf("method (%s) != expectedMethod (%s)", r.Method, expectedMethod)
	}
}

func ensureHeader(t *testing.T, r *http.Request, header, expectedValue string) {
	t.Helper()
	if r.Header.Get(header) != expectedValue {
		t.Fatalf("token (%s) != expectedValue (%s)", r.Header.Get("token"), expectedValue)
	}
}

func ensurePayload(t *testing.T, r *http.Request, expect []byte) {
	t.Helper()
	data, errRead := ioutil.ReadAll(r.Body)
	if errRead != nil {
		t.Fatal(errRead)
	}

	if !bytes.Equal(data, expect) {
		t.Fatalf("payload doesn't match\nhas: >%s<\nexpect: >%s<", data, expect)
	}
}

func ensureURLParameter(t *testing.T, r *http.Request, parameter, expectedValue string) {
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
	}

	t.Fatalf("error was not of type BaseResponse, error string is: %s", err)
}

func getBaseResponse(wantedError ErrorCode) string {
	b, _ := json.Marshal(BaseResponse{
		Response: wantedError.InternalCode(),
		Code:     wantedError,
		ServerID: "serverID1",
		Message:  wantedError.String(),
		DateTime: time.Now(),
	})

	return string(b)
}

func TestGetTokenOK(t *testing.T) {
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
		t.FailNow()
	}

}

func TestGetTokenInvalidUser(t *testing.T) {
	setup()

	mux.HandleFunc(fmt.Sprint("/", APIVersion, "/token"),
		func(w http.ResponseWriter, r *http.Request) {
			ensureMethod(t, r, "POST")

			var tokenResp TokenResponse

			errDecode := json.NewDecoder(r.Body).Decode(&tokenResp)
			if errDecode != nil {
				t.Fatal(errDecode)
			}

			fmt.Fprint(w, getBaseResponse(ErrInvalidUser))
		},
	)

	_, errToken := client.GetToken("user1", "pass1")
	ensureError(t, errToken, ErrInvalidUser)
}

func TestGetStatusOK(t *testing.T) {
	setup()

	mux.HandleFunc(fmt.Sprint("/", APIVersion, "/status"),
		func(w http.ResponseWriter, r *http.Request) {
			ensureMethod(t, r, "GET")
			ensureHeader(t, r, "token", "d97c908ed44c25fdca302612c70584c8d5acd47a")

			fmt.Fprint(w, `{"account":{"expires":"2014-09-26T19:07:28Z","messages":[],"maxLineups":4,"nextSuggestedConnectTime":"2014-07-29T22:43:22Z"},"lineups":[],"lastDataUpdate":"2014-07-28T14:48:59Z","notifications":[],"systemStatus":[{"date":"2012-12-17T16:24:47Z","status":"Online","details":"All servers running normally."}],"serverID":"serverID1","code":0}`)
		},
	)

	status, err := client.GetStatus()
	if err != nil {
		t.Fatal(err)
	}

	if len(status.SystemStatus) != 1 {
		t.FailNow()
	} else if status.SystemStatus[0].Details != "All servers running normally." {
		t.FailNow()
	}
}
