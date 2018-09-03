package schedulesdirect

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Date is a Schedules Direct specific date format (YYYY[-MM-DD]) with (Un)MarshalJSON functions.
type Date struct {
	*time.Time
	fmt string
}

// MarshalJSON formats the underlying time.Time to Schedule Direct's format.
func (p Date) MarshalJSON() ([]byte, error) {
	str := "\"" + p.Format(p.fmt) + "\""

	return []byte(str), nil
}

// UnmarshalJSON converts Schedule Direct's format to a time.Time.
func (p *Date) UnmarshalJSON(text []byte) error {
	dateFormat := "2006-01-02"

	str, unquoteErr := strconv.Unquote(string(text))
	if unquoteErr != nil {
		return unquoteErr
	}

	if len(str) == 4 {
		dateFormat = "2006"
	}

	v, e := time.Parse(dateFormat, str)
	if e != nil {
		return fmt.Errorf("schedulesdirect.Date should be a time, error value is: %s", text)
	}

	*p = Date{&v, dateFormat}

	return nil
}

// jsonInt is a int64 which unmarshals from JSON
// as either unquoted or quoted (with any amount
// of internal leading/trailing whitespace).
// Originally found at https://bit.ly/2NkJ0SK and
// https://play.golang.org/p/KNPxDL1yqL
type jsonInt int64

func (f jsonInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(int64(f))
}

func (f *jsonInt) UnmarshalJSON(data []byte) error {
	var v int64

	data = bytes.Trim(data, `" `)

	err := json.Unmarshal(data, &v)
	*f = jsonInt(v)
	return err
}

// ConvertibleBoolean is a helper type to allow JSON documents using 0/1, "true" and "false" or "yes" and "no" be converted to bool.
type ConvertibleBoolean struct {
	bool
	quoted bool
}

// MarshalJSON returns a 0 or 1 depending on bool state.
func (bit ConvertibleBoolean) MarshalJSON() ([]byte, error) {
	var bitSetVar int8
	if bit.bool {
		bitSetVar = 1
	}

	if bit.quoted {
		return json.Marshal(fmt.Sprint(bitSetVar))
	}

	return json.Marshal(bitSetVar)
}

// UnmarshalJSON converts a 0, 1, true or false into a bool
func (bit *ConvertibleBoolean) UnmarshalJSON(data []byte) error {
	bit.quoted = strings.Contains(string(data), `"`)
	// Bools are sometimes quoted, sometimes not, lets just always remove quotes just in case...
	asString := strings.Replace(string(data), `"`, "", -1)
	if asString == "1" || asString == "true" || asString == "yes" {
		bit.bool = true
	} else if asString == "0" || asString == "false" || asString == "no" {
		bit.bool = false
	} else {
		return fmt.Errorf("Boolean unmarshal error: invalid input %s", asString)
	}
	return nil
}

// BaseResponse contains the fields that every request is expected to return.
type BaseResponse struct {
	Response string    `json:"response,omitempty"`
	Code     ErrorCode `json:"code,omitempty"`
	ServerID string    `json:"serverID,omitempty"`
	Message  string    `json:"message,omitempty"`
	DateTime time.Time `json:"datetime,omitempty"`
}

// Error returns a error string.
func (e BaseResponse) Error() string {
	return e.Code.Error()
}

// A TokenResponse stores the response for token request.
type TokenResponse struct {
	*BaseResponse

	Token string `json:"token,omitempty"`
}

// A VersionResponse stores the response for a version request.
type VersionResponse struct {
	*BaseResponse

	Client  string `json:"client,omitempty"`
	Version string `json:"version,omitempty"`
}

// A StatusResponse stores the message after requesting system
// status.  SystemStatus[0].Status should be "Online" before proceeding.
type StatusResponse struct {
	*BaseResponse

	Account        *AccountInfo `json:"account,omitempty"`
	Lineups        []Lineup     `json:"lineups,omitempty"`
	LastDataUpdate time.Time    `json:"lastDataUpdate,omitempty"`
	Notifications  []string     `json:"notifications,omitempty"`
	SystemStatus   []Status     `json:"systemStatus,omitempty"`
}

// A Status stores the message system status information
// usually as part of a StatusResponse.
type Status struct {
	Date    *time.Time `json:"date,omitempty"`
	Status  string     `json:"status,omitempty"`
	Details string     `json:"details,omitempty"`
}

// An AccountInfo stores the message account information
// usually as part of a StatusResponse.
type AccountInfo struct {
	Expires    string   `json:"expires,omitempty"`
	Messages   []string `json:"messages,omitempty"`
	MaxLineups int      `json:"maxLineups,omitempty"`
}
