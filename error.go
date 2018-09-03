package schedulesdirect

import (
	"fmt"
	"strconv"
)

// A ErrorCode is an unsigned 32-bit error code as defined in the Schedules Direct JSON specification.
//
// This implementation was inspired by gRPC's codes package: https://github.com/grpc/grpc-go/tree/master/codes.
type ErrorCode uint32

const (

	// ErrOK is returned when everything is OK.
	ErrOK ErrorCode = 0

	// ErrInvalidJSON is returned when unable to decode JSON.
	ErrInvalidJSON ErrorCode = 1001

	// ErrDeflateRequired is returned when did not receive Accept-Encoding: deflate in request.
	ErrDeflateRequired ErrorCode = 1002

	// ErrTokenMissing is returned when token required but not provided in request header.
	ErrTokenMissing ErrorCode = 1004

	// ErrUnsupportedCommand is returned when unsupported command.
	ErrUnsupportedCommand ErrorCode = 2000

	// ErrRequiredActionMissing is returned when request is missing an action to take.
	ErrRequiredActionMissing ErrorCode = 2001

	// ErrRequiredRequestMissing is returned when did not receive request.
	ErrRequiredRequestMissing ErrorCode = 2002

	// ErrRequiredParameterMissingCountry is returned when in order to search for lineups, you must supply a 3-letter country parameter.
	ErrRequiredParameterMissingCountry ErrorCode = 2004

	// ErrRequiredParameterMissingPostalCode is returned when in order to search for lineups, you must supply a postal code parameter.
	ErrRequiredParameterMissingPostalCode ErrorCode = 2005

	// ErrRequiredParameterMissingMessageID is returned when in order to delete a message you must supply the messageID.
	ErrRequiredParameterMissingMessageID ErrorCode = 2006

	// ErrInvalidParameterCountry is returned when the COUNTRY parameter was not a ISO-3166-1 alpha 3 code.
	// See http://en.wikipedia.org/wiki/ISO_3166-1_alpha-3 for more information.
	ErrInvalidParameterCountry ErrorCode = 2050

	// ErrInvalidParameterPostalCode is returned when the POSTALCODE parameter must was invalid for the country searched.
	// Post message to http://forums.schedulesdirect.org/viewforum.php?f=6 if you are having issues.
	ErrInvalidParameterPostalCode ErrorCode = 2051

	// ErrInvalidParameterFetchType is returned when you didn't provide a fetchtype Schedules Direct knows how to handle.
	ErrInvalidParameterFetchType ErrorCode = 2052

	// ErrDuplicateLineup is returned when a lineup already existed in account.
	ErrDuplicateLineup ErrorCode = 2100

	// ErrLineupNotFound is returned when lineup is not in an account.
	// Add lineup to account before requesting mapping.
	ErrLineupNotFound ErrorCode = 2101

	// ErrUnknownLineup is returned when invalid lineup requested.
	// Check your COUNTRY / POSTALCODE combination for validity.
	ErrUnknownLineup ErrorCode = 2102

	// ErrInvalidLineupDelete is returned when delete of lineup not in account.
	ErrInvalidLineupDelete ErrorCode = 2103

	// ErrLineupWrongFormat is returned when lineup must be formatted COUNTRY-LINEUP-DEVICE or COUNTRY-OTA-POSTALCODE.
	ErrLineupWrongFormat ErrorCode = 2104

	// ErrInvalidLineup is returned when the lineup you submitted doesn't exist.
	ErrInvalidLineup ErrorCode = 2105

	// ErrLineupDeleted is returned when the lineup you requested has been deleted from the server.
	ErrLineupDeleted ErrorCode = 2106

	// ErrLineupQueued is returned when the lineup is being generated on the server.
	// Please retry.
	ErrLineupQueued ErrorCode = 2107

	// ErrInvalidCountry is returned when the country you requested is either mis-typed or does not have valid data.
	ErrInvalidCountry ErrorCode = 2108

	// ErrStationIDNotFound is returned when the stationID you requested is not in any of your lineups.
	ErrStationIDNotFound ErrorCode = 2200

	// ErrServiceOffline is returned when server offline for maintenance.
	ErrServiceOffline ErrorCode = 3000

	// ErrAccountExpired is returned when account expired.
	ErrAccountExpired ErrorCode = 4001

	// ErrInvalidHash is returned when password hash must be lowercase 40 character sha1_hex of password.
	ErrInvalidHash ErrorCode = 4002

	// ErrInvalidUser is returned when invalid username or password.
	ErrInvalidUser ErrorCode = 4003

	// ErrAccountLockout is returned when too many login failures.
	// Locked for 15 minutes.
	ErrAccountLockout ErrorCode = 4004

	// ErrAccountDisabled is returned when account has been disabled.
	// Please contact Schedules Direct support: admin@schedulesdirect.org for more information.
	ErrAccountDisabled ErrorCode = 4005

	// ErrTokenExpired is returned when token has expired.
	// Request new token.
	ErrTokenExpired ErrorCode = 4006

	// ErrMaxLineupChangesReached is returned when exceeded maximum number of lineup changes for today.
	ErrMaxLineupChangesReached ErrorCode = 4100

	// ErrMaxLineups is returned when exceeded number of lineups for this account.
	ErrMaxLineups ErrorCode = 4101

	// ErrNoLineups is returned when no lineups have been added to this account.
	ErrNoLineups ErrorCode = 4102

	// ErrImageNotFound is returned when could not find requested image.
	// Post message to http://forums.schedulesdirect.org/viewforum.php?f=6 if you are having issues.
	ErrImageNotFound ErrorCode = 5000

	// ErrInvalidProgramID is returned when could not find requested programID.
	// Permanent  failure.
	ErrInvalidProgramID ErrorCode = 6000

	// ErrProgramIDQueued is returned when programID should exist at the server, but doesn't.
	// The server will regenerate the JSON for the program, so your application should retry.
	ErrProgramIDQueued ErrorCode = 6001

	// ErrScheduleNotFound is returned when the schedule you requested should be available.
	// Post message to http://forums.schedulesdirect.org/viewforum.php?f=6
	ErrScheduleNotFound ErrorCode = 7000

	// ErrInvalidScheduleRequest is returned when the server can't determine whether your schedule is valid or not.
	// Open a support ticket.
	ErrInvalidScheduleRequest ErrorCode = 7010

	// ErrScheduleRangeExceeded is returned when the date that you've requested is outside of the range of the data for that stationID.
	ErrScheduleRangeExceeded ErrorCode = 7020

	// ErrScheduleNotInLineup is returned when you have requested a schedule which is not in any of your configured lineups.
	ErrScheduleNotInLineup ErrorCode = 7030

	// ErrScheduleQueued is returned when the schedule you requested has been queued for generation but is not yet ready for download.
	// Retry.
	ErrScheduleQueued ErrorCode = 7100

	// ErrHCF is returned when Unknown error.
	// Open support ticket.
	ErrHCF ErrorCode = 9999

	_maxCode = 10000
)

// UnmarshalJSON unmarshals b into the ErrorCode.
func (c *ErrorCode) UnmarshalJSON(b []byte) error {
	// From json.Unmarshaler: By convention, to approximate the behavior of
	// Unmarshal itself, Unmarshalers implement UnmarshalJSON([]byte("null")) as
	// a no-op.
	if string(b) == "null" {
		return nil
	}
	if c == nil {
		return fmt.Errorf("nil receiver passed to UnmarshalJSON")
	}

	if ci, err := strconv.ParseUint(string(b), 10, 32); err == nil {
		if ci >= _maxCode {
			return fmt.Errorf("invalid code: %q", ci)
		}

		*c = ErrorCode(ci)
		return nil
	}

	if jc, ok := strToCode[string(b)]; ok {
		*c = jc
		return nil
	}
	return fmt.Errorf("invalid code: %q", string(b))
}

var intToCode = map[int32]ErrorCode{
	0:    ErrOK,
	1001: ErrInvalidJSON,
	1002: ErrDeflateRequired,
	1004: ErrTokenMissing,
	2000: ErrUnsupportedCommand,
	2001: ErrRequiredActionMissing,
	2002: ErrRequiredRequestMissing,
	2004: ErrRequiredParameterMissingCountry,
	2005: ErrRequiredParameterMissingPostalCode,
	2006: ErrRequiredParameterMissingMessageID,
	2050: ErrInvalidParameterCountry,
	2051: ErrInvalidParameterPostalCode,
	2052: ErrInvalidParameterFetchType,
	2100: ErrDuplicateLineup,
	2101: ErrLineupNotFound,
	2102: ErrUnknownLineup,
	2103: ErrInvalidLineupDelete,
	2104: ErrLineupWrongFormat,
	2105: ErrInvalidLineup,
	2106: ErrLineupDeleted,
	2107: ErrLineupQueued,
	2108: ErrInvalidCountry,
	2200: ErrStationIDNotFound,
	3000: ErrServiceOffline,
	4001: ErrAccountExpired,
	4002: ErrInvalidHash,
	4003: ErrInvalidUser,
	4004: ErrAccountLockout,
	4005: ErrAccountDisabled,
	4006: ErrTokenExpired,
	4100: ErrMaxLineupChangesReached,
	4101: ErrMaxLineups,
	4102: ErrNoLineups,
	5000: ErrImageNotFound,
	6000: ErrInvalidProgramID,
	6001: ErrProgramIDQueued,
	7000: ErrScheduleNotFound,
	7010: ErrInvalidScheduleRequest,
	7020: ErrScheduleRangeExceeded,
	7030: ErrScheduleNotInLineup,
	7100: ErrScheduleQueued,
	9999: ErrHCF,
}

var internalToInt = map[string]int32{
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

var strToCode = map[string]ErrorCode{
	`"OK"`:                                    ErrOK,
	`"INVALID_JSON"`:                          ErrInvalidJSON,
	`"DEFLATE_REQUIRED"`:                      ErrDeflateRequired,
	`"TOKEN_MISSING"`:                         ErrTokenMissing,
	`"UNSUPPORTED_COMMAND"`:                   ErrUnsupportedCommand,
	`"REQUIRED_ACTION_MISSING"`:               ErrRequiredActionMissing,
	`"REQUIRED_REQUEST_MISSING"`:              ErrRequiredRequestMissing,
	`"REQUIRED_PARAMETER_MISSING:COUNTRY"`:    ErrRequiredParameterMissingCountry,
	`"REQUIRED_PARAMETER_MISSING:POSTALCODE"`: ErrRequiredParameterMissingPostalCode,
	`"REQUIRED_PARAMETER_MISSING:MSGID"`:      ErrRequiredParameterMissingMessageID,
	`"INVALID_PARAMETER:COUNTRY"`:             ErrInvalidParameterCountry,
	`"INVALID_PARAMETER:POSTALCODE"`:          ErrInvalidParameterPostalCode,
	`"INVALID_PARAMETER:FETCHTYPE"`:           ErrInvalidParameterFetchType,
	`"DUPLICATE_LINEUP"`:                      ErrDuplicateLineup,
	`"LINEUP_NOT_FOUND"`:                      ErrLineupNotFound,
	`"UNKNOWN_LINEUP"`:                        ErrUnknownLineup,
	`"INVALID_LINEUP_DELETE"`:                 ErrInvalidLineupDelete,
	`"LINEUP_WRONG_FORMAT"`:                   ErrLineupWrongFormat,
	`"INVALID_LINEUP"`:                        ErrInvalidLineup,
	`"LINEUP_DELETED"`:                        ErrLineupDeleted,
	`"LINEUP_QUEUED"`:                         ErrLineupQueued,
	`"INVALID_COUNTRY"`:                       ErrInvalidCountry,
	`"STATIONID_NOT_FOUND"`:                   ErrStationIDNotFound,
	`"SERVICE_OFFLINE"`:                       ErrServiceOffline,
	`"ACCOUNT_EXPIRED"`:                       ErrAccountExpired,
	`"INVALID_HASH"`:                          ErrInvalidHash,
	`"INVALID_USER"`:                          ErrInvalidUser,
	`"ACCOUNT_LOCKOUT"`:                       ErrAccountLockout,
	`"ACCOUNT_DISABLED"`:                      ErrAccountDisabled,
	`"TOKEN_EXPIRED"`:                         ErrTokenExpired,
	`"MAX_LINEUP_CHANGES_REACHED"`:            ErrMaxLineupChangesReached,
	`"MAX_LINEUPS"`:                           ErrMaxLineups,
	`"NO_LINEUPS"`:                            ErrNoLineups,
	`"IMAGE_NOT_FOUND"`:                       ErrImageNotFound,
	`"INVALID_PROGRAMID"`:                     ErrInvalidProgramID,
	`"PROGRAMID_QUEUED"`:                      ErrProgramIDQueued,
	`"SCHEDULE_NOT_FOUND"`:                    ErrScheduleNotFound,
	`"INVALID_SCHEDULE_REQUEST"`:              ErrInvalidScheduleRequest,
	`"SCHEDULE_RANGE_EXCEEDED"`:               ErrScheduleRangeExceeded,
	`"SCHEDULE_NOT_IN_LINEUP"`:                ErrScheduleNotInLineup,
	`"SCHEDULE_QUEUED"`:                       ErrScheduleQueued,
	`"HCF"`:                                   ErrHCF,
}

var codeToInternal = map[ErrorCode]string{
	ErrOK:                                 "OK",
	ErrInvalidJSON:                        "INVALID_JSON",
	ErrDeflateRequired:                    "DEFLATE_REQUIRED",
	ErrTokenMissing:                       "TOKEN_MISSING",
	ErrUnsupportedCommand:                 "UNSUPPORTED_COMMAND",
	ErrRequiredActionMissing:              "REQUIRED_ACTION_MISSING",
	ErrRequiredRequestMissing:             "REQUIRED_REQUEST_MISSING",
	ErrRequiredParameterMissingCountry:    "REQUIRED_PARAMETER_MISSING:COUNTRY",
	ErrRequiredParameterMissingPostalCode: "REQUIRED_PARAMETER_MISSING:POSTALCODE",
	ErrRequiredParameterMissingMessageID:  "REQUIRED_PARAMETER_MISSING:MSGID",
	ErrInvalidParameterCountry:            "INVALID_PARAMETER:COUNTRY",
	ErrInvalidParameterPostalCode:         "INVALID_PARAMETER:POSTALCODE",
	ErrInvalidParameterFetchType:          "INVALID_PARAMETER:FETCHTYPE",
	ErrDuplicateLineup:                    "DUPLICATE_LINEUP",
	ErrLineupNotFound:                     "LINEUP_NOT_FOUND",
	ErrUnknownLineup:                      "UNKNOWN_LINEUP",
	ErrInvalidLineupDelete:                "INVALID_LINEUP_DELETE",
	ErrLineupWrongFormat:                  "LINEUP_WRONG_FORMAT",
	ErrInvalidLineup:                      "INVALID_LINEUP",
	ErrLineupDeleted:                      "LINEUP_DELETED",
	ErrLineupQueued:                       "LINEUP_QUEUED",
	ErrInvalidCountry:                     "INVALID_COUNTRY",
	ErrStationIDNotFound:                  "STATIONID_NOT_FOUND",
	ErrServiceOffline:                     "SERVICE_OFFLINE",
	ErrAccountExpired:                     "ACCOUNT_EXPIRED",
	ErrInvalidHash:                        "INVALID_HASH",
	ErrInvalidUser:                        "INVALID_USER",
	ErrAccountLockout:                     "ACCOUNT_LOCKOUT",
	ErrAccountDisabled:                    "ACCOUNT_DISABLED",
	ErrTokenExpired:                       "TOKEN_EXPIRED",
	ErrMaxLineupChangesReached:            "MAX_LINEUP_CHANGES_REACHED",
	ErrMaxLineups:                         "MAX_LINEUPS",
	ErrNoLineups:                          "NO_LINEUPS",
	ErrImageNotFound:                      "IMAGE_NOT_FOUND",
	ErrInvalidProgramID:                   "INVALID_PROGRAMID",
	ErrProgramIDQueued:                    "PROGRAMID_QUEUED",
	ErrScheduleNotFound:                   "SCHEDULE_NOT_FOUND",
	ErrInvalidScheduleRequest:             "INVALID_SCHEDULE_REQUEST",
	ErrScheduleRangeExceeded:              "SCHEDULE_RANGE_EXCEEDED",
	ErrScheduleNotInLineup:                "SCHEDULE_NOT_IN_LINEUP",
	ErrScheduleQueued:                     "SCHEDULE_QUEUED",
	ErrHCF:                                "HCF",
}

var codeToMessage = map[ErrorCode]string{
	ErrOK:                                 "OK",
	ErrInvalidJSON:                        "Unable to decode JSON",
	ErrDeflateRequired:                    "Did not receive Accept-Encoding: deflate in request.",
	ErrTokenMissing:                       "Token required but not provided in request header.",
	ErrUnsupportedCommand:                 "Unsupported command",
	ErrRequiredActionMissing:              "Request is missing an action to take.",
	ErrRequiredRequestMissing:             "Did not receive request.",
	ErrRequiredParameterMissingCountry:    "In order to search for lineups, you must supply a 3-letter country parameter.",
	ErrRequiredParameterMissingPostalCode: "In order to search for lineups, you must supply a postal code parameter.",
	ErrRequiredParameterMissingMessageID:  "In order to delete a message you must supply the messageID.",
	ErrInvalidParameterCountry:            "The COUNTRY parameter must be ISO-3166-1 alpha 3. See http://en.wikipedia.org/wiki/ISO_3166-1_alpha-3",
	ErrInvalidParameterPostalCode:         "The POSTALCODE parameter must be valid for the country you are searching. Post message to http://forums.schedulesdirect.org/viewforum.php?f=6 if you are having issues.",
	ErrInvalidParameterFetchType:          "You didn't provide a fetchtype I know how to handle.",
	ErrDuplicateLineup:                    "Lineup already in account.",
	ErrLineupNotFound:                     "Lineup not in account. Add lineup to account before requesting mapping.",
	ErrUnknownLineup:                      "Invalid lineup requested. Check your COUNTRY / POSTALCODE combination for validity.",
	ErrInvalidLineupDelete:                "Delete of lineup not in account.",
	ErrLineupWrongFormat:                  "Lineup must be formatted COUNTRY-LINEUP-DEVICE or COUNTRY-OTA-POSTALCODE",
	ErrInvalidLineup:                      "The lineup you submitted doesn't exist.",
	ErrLineupDeleted:                      "The lineup you requested has been deleted from the server.",
	ErrLineupQueued:                       "The lineup is being generated on the server. Please retry.",
	ErrInvalidCountry:                     "The country you requested is either mis-typed or does not have valid data.",
	ErrStationIDNotFound:                  "The stationID you requested is not in any of your lineups.",
	ErrServiceOffline:                     "Server offline for maintenance.",
	ErrAccountExpired:                     "Account expired.",
	ErrInvalidHash:                        "Password hash must be lowercase 40 character sha1_hex of password.",
	ErrInvalidUser:                        "Invalid username or password.",
	ErrAccountLockout:                     "Too many login failures. Locked for 15 minutes.",
	ErrAccountDisabled:                    "Account has been disabled. Please contact Schedules Direct support: admin@schedulesdirect.org for more information.",
	ErrTokenExpired:                       "Token has expired. Request new token.",
	ErrMaxLineupChangesReached:            "Exceeded maximum number of lineup changes for today.",
	ErrMaxLineups:                         "Exceeded number of lineups for this account.",
	ErrNoLineups:                          "No lineups have been added to this account.",
	ErrImageNotFound:                      "Could not find requested image. Post message to http://forums.schedulesdirect.org/viewforum.php?f=6 if you are having issues.",
	ErrInvalidProgramID:                   "Could not find requested programID. Permanent failure.",
	ErrProgramIDQueued:                    "ProgramID should exist at the server, but doesn't. The server will regenerate the JSON for the program, so your application should retry.",
	ErrScheduleNotFound:                   "The schedule you requested should be available. Post message to http://forums.schedulesdirect.org/viewforum.php?f=6",
	ErrInvalidScheduleRequest:             "The server can't determine whether your schedule is valid or not. Open a support ticket.",
	ErrScheduleRangeExceeded:              "The date that you've requested is outside of the range of the data for that stationID.",
	ErrScheduleNotInLineup:                "You have requested a schedule which is not in any of your configured lineups.",
	ErrScheduleQueued:                     "The schedule you requested has been queued for generation but is not yet ready for download. Retry.",
	ErrHCF:                                "Unknown error. Open support ticket.",
}

// InternalCode returns the Schedules Direct internal error code, like "INVALID_JSON".
func (c ErrorCode) InternalCode() string {
	if val, ok := codeToInternal[c]; ok {
		return val
	}
	return "Unknown ErrorCode(" + strconv.FormatInt(int64(c), 10) + ")"
}

func (c ErrorCode) String() string {
	if val, ok := codeToMessage[c]; ok {
		return val
	}
	return "Unknown ErrorCode(" + strconv.FormatInt(int64(c), 10) + ")"
}

func (c ErrorCode) Error() string {
	return fmt.Sprintf("%s (message: %s, code: %d)", c.String(), c.InternalCode(), c)
}
