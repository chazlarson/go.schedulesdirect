package schedulesdirect

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestErrorUnmarshalJSON(t *testing.T) {
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
