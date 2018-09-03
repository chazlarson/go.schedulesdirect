package schedulesdirect

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestErrorUnmarshalJSON(t *testing.T) {
	internalToInt := map[string]int32{
		"OK":                                    0,
		"INVALID_JSON":                          1001,
		"DEFLATE_REQUIRED":                      1002,
		"TOKEN_MISSING":                         1004,
		"UNSUPPORTED_COMMAND":                   2000,
		"REQUIRED_ACTION_MISSING":               2001,
		"REQUIRED_REQUEST_MISSING":              2002,
		"REQUIRED_PARAMETER_MISSING:COUNTRY":    2004,
		"REQUIRED_PARAMETER_MISSING:POSTALCODE": 2005,
		"REQUIRED_PARAMETER_MISSING:MSGID":      2006,
		"INVALID_PARAMETER:COUNTRY":             2050,
		"INVALID_PARAMETER:POSTALCODE":          2051,
		"INVALID_PARAMETER:FETCHTYPE":           2052,
		"DUPLICATE_LINEUP":                      2100,
		"LINEUP_NOT_FOUND":                      2101,
		"UNKNOWN_LINEUP":                        2102,
		"INVALID_LINEUP_DELETE":                 2103,
		"LINEUP_WRONG_FORMAT":                   2104,
		"INVALID_LINEUP":                        2105,
		"LINEUP_DELETED":                        2106,
		"LINEUP_QUEUED":                         2107,
		"INVALID_COUNTRY":                       2108,
		"STATIONID_NOT_FOUND":                   2200,
		"SERVICE_OFFLINE":                       3000,
		"ACCOUNT_EXPIRED":                       4001,
		"INVALID_HASH":                          4002,
		"INVALID_USER":                          4003,
		"ACCOUNT_LOCKOUT":                       4004,
		"ACCOUNT_DISABLED":                      4005,
		"TOKEN_EXPIRED":                         4006,
		"MAX_LINEUP_CHANGES_REACHED":            4100,
		"MAX_LINEUPS":                           4101,
		"NO_LINEUPS":                            4102,
		"IMAGE_NOT_FOUND":                       5000,
		"INVALID_PROGRAMID":                     6000,
		"PROGRAMID_QUEUED":                      6001,
		"SCHEDULE_NOT_FOUND":                    7000,
		"INVALID_SCHEDULE_REQUEST":              7010,
		"SCHEDULE_RANGE_EXCEEDED":               7020,
		"SCHEDULE_NOT_IN_LINEUP":                7030,
		"SCHEDULE_QUEUED":                       7100,
		"HCF":                                   9999,
	}

	for s, v := range internalToInt {
		want := ErrorCode(v)
		var got ErrorCode
		if err := got.UnmarshalJSON([]byte(`"` + s + `"`)); err != nil || got != want {
			t.Errorf("got.UnmarshalJSON(%q) = %v; want <nil>.  got=%v; want %v", s, err, got, want)
		}
	}
}

func TestErrorJSONUnmarshal(t *testing.T) {
	var got []ErrorCode
	want := []ErrorCode{ErrOK, ErrLineupNotFound, ErrHCF, ErrUnsupportedCommand}
	in := `["OK", "LINEUP_NOT_FOUND", "HCF", "UNSUPPORTED_COMMAND"]`
	err := json.Unmarshal([]byte(in), &got)
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("json.Unmarshal(%q, &got) = %v; want <nil>.  got=%v; want %v", in, err, got, want)
	}
}

func TestErrorUnmarshalJSON_NilReceiver(t *testing.T) {
	var got *ErrorCode
	in := ErrOK.String()
	if err := got.UnmarshalJSON([]byte(in)); err == nil {
		t.Errorf("got.UnmarshalJSON(%q) = nil; want <non-nil>.  got=%v", in, got)
	}
}

func TestErrorUnmarshalJSON_UnknownInput(t *testing.T) {
	var got ErrorCode
	for _, in := range [][]byte{[]byte(""), []byte("xxx"), []byte("ErrorCode(17)"), nil} {
		if err := got.UnmarshalJSON([]byte(in)); err == nil {
			t.Errorf("got.UnmarshalJSON(%q) = nil; want <non-nil>.  got=%v", in, got)
		}
	}
}

func TestErrorUnmarshalJSON_MarshalUnmarshal(t *testing.T) {
	for i := 0; i < _maxCode; i++ {
		var cUnMarshaled ErrorCode
		c := ErrorCode(i)

		cJSON, err := json.Marshal(c)
		if err != nil {
			t.Errorf("marshalling %q failed: %v", c, err)
		}

		if err := json.Unmarshal(cJSON, &cUnMarshaled); err != nil {
			t.Errorf("unmarshalling ErrorCode failed: %s", err)
		}

		if c != cUnMarshaled {
			t.Errorf("ErrorCode is %q after marshalling/unmarshalling, expected %q", cUnMarshaled, c)
		}
	}
}
