package schedulesdirect

import "strconv"

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

	// InvalidParameterCountry is returned when the COUNTRY parameter was not a ISO-3166-1 alpha 3 code.
	// ErrInvalidParameterCountry http://en.wikipedia.org/wiki/ISO_3166-1_alpha-3 for more information.
	ErrInvalidParameterCountry ErrorCode = 2050

	// InvalidParameterPostalCode is returned when the POSTALCODE parameter must was invalid for the country searched.
	// ErrInvalidParameterPostalCode message to http://forums.schedulesdirect.org/viewforum.php?f=6 if you are having issues.
	ErrInvalidParameterPostalCode ErrorCode = 2051

	// ErrInvalidParameterFetchType is returned when you didn't provide a fetchtype Schedules Direct knows how to handle.
	ErrInvalidParameterFetchType ErrorCode = 2052

	// ErrDuplicateLineup is returned when a lineup already existed in account.
	ErrDuplicateLineup ErrorCode = 2100

	// LineupNotFound is returned when lineup is not in an account.
	// ErrLineupNotFound lineup to account before requesting mapping.
	ErrLineupNotFound ErrorCode = 2101

	// UnknownLineup is returned when invalid lineup requested.
	// ErrUnknownLineup your COUNTRY / POSTALCODE combination for validity.
	ErrUnknownLineup ErrorCode = 2102

	// ErrInvalidLineupDelete is returned when delete of lineup not in account.
	ErrInvalidLineupDelete ErrorCode = 2103

	// ErrLineupWrongFormat is returned when lineup must be formatted COUNTRY-LINEUP-DEVICE or COUNTRY-OTA-POSTALCODE.
	ErrLineupWrongFormat ErrorCode = 2104

	// ErrInvalidLineup is returned when the lineup you submitted doesn't exist.
	ErrInvalidLineup ErrorCode = 2105

	// ErrLineupDeleted is returned when the lineup you requested has been deleted from the server.
	ErrLineupDeleted ErrorCode = 2106

	// LineupQueued is returned when the lineup is being generated on the server.
	// ErrLineupQueued retry.
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

	// AccountLockout is returned when too many login failures.
	// ErrAccountLockout for 15 minutes.
	ErrAccountLockout ErrorCode = 4004

	// AccountDisabled is returned when account has been disabled.
	// ErrAccountDisabled contact Schedules Direct support: admin@schedulesdirect.org for more information.
	ErrAccountDisabled ErrorCode = 4005

	// TokenExpired is returned when token has expired.
	// ErrTokenExpired new token.
	ErrTokenExpired ErrorCode = 4006

	// ErrMaxLineupChangesReached is returned when exceeded maximum number of lineup changes for today.
	ErrMaxLineupChangesReached ErrorCode = 4100

	// ErrMaxLineups is returned when exceeded number of lineups for this account.
	ErrMaxLineups ErrorCode = 4101

	// ErrNoLineups is returned when no lineups have been added to this account.
	ErrNoLineups ErrorCode = 4102

	// ImageNotFound is returned when could not find requested image.
	// ErrImageNotFound message to http://forums.schedulesdirect.org/viewforum.php?f=6 if you are having issues.
	ErrImageNotFound ErrorCode = 5000

	// InvalidProgramID is returned when could not find requested programID.
	// ErrInvalidProgramID failure.
	ErrInvalidProgramID ErrorCode = 6000

	// ProgramidQueued is returned when programID should exist at the server, but doesn't.
	// ErrProgramidQueued server will regenerate the JSON for the program, so your application should retry.
	ErrProgramidQueued ErrorCode = 6001

	// ScheduleNotFound is returned when the schedule you requested should be available.
	// ErrScheduleNotFound message to http://forums.schedulesdirect.org/viewforum.php?f=6
	ErrScheduleNotFound ErrorCode = 7000

	// InvalidScheduleRequest is returned when the server can't determine whether your schedule is valid or not.
	// ErrInvalidScheduleRequest a support ticket.
	ErrInvalidScheduleRequest ErrorCode = 7010

	// ErrScheduleRangeExceeded is returned when the date that you've requested is outside of the range of the data for that stationID.
	ErrScheduleRangeExceeded ErrorCode = 7020

	// ErrScheduleNotInLineup is returned when you have requested a schedule which is not in any of your configured lineups.
	ErrScheduleNotInLineup ErrorCode = 7030

	// ScheduleQueued is returned when the schedule you requested has been queued for generation but is not yet ready for download.
	// ErrScheduleQueued.
	ErrScheduleQueued ErrorCode = 7100

	// HCF is returned when Unknown error.
	// ErrHCF support ticket.
	ErrHCF ErrorCode = 9999
)

// InternalCode returns the Schedules Direct internal error code, like "INVALID_JSON".
func (c ErrorCode) InternalCode() string {
	switch c {
	case ErrOK:
		return "OK"
	case ErrInvalidJSON:
		return "INVALID_JSON"
	case ErrDeflateRequired:
		return "DEFLATE_REQUIRED"
	case ErrTokenMissing:
		return "TOKEN_MISSING"
	case ErrUnsupportedCommand:
		return "UNSUPPORTED_COMMAND"
	case ErrRequiredActionMissing:
		return "REQUIRED_ACTION_MISSING"
	case ErrRequiredRequestMissing:
		return "REQUIRED_REQUEST_MISSING"
	case ErrRequiredParameterMissingCountry:
		return "REQUIRED_PARAMETER_MISSING"
	case ErrRequiredParameterMissingPostalCode:
		return "REQUIRED_PARAMETER_MISSING"
	case ErrRequiredParameterMissingMessageID:
		return "REQUIRED_PARAMETER_MISSING"
	case ErrInvalidParameterCountry:
		return "INVALID_PARAMETER"
	case ErrInvalidParameterPostalCode:
		return "INVALID_PARAMETER"
	case ErrInvalidParameterFetchType:
		return "INVALID_PARAMETER"
	case ErrDuplicateLineup:
		return "DUPLICATE_LINEUP"
	case ErrLineupNotFound:
		return "LINEUP_NOT_FOUND"
	case ErrUnknownLineup:
		return "UNKNOWN_LINEUP"
	case ErrInvalidLineupDelete:
		return "INVALID_LINEUP_DELETE"
	case ErrLineupWrongFormat:
		return "LINEUP_WRONG_FORMAT"
	case ErrInvalidLineup:
		return "INVALID_LINEUP"
	case ErrLineupDeleted:
		return "LINEUP_DELETED"
	case ErrLineupQueued:
		return "LINEUP_QUEUED"
	case ErrInvalidCountry:
		return "INVALID_COUNTRY"
	case ErrStationIDNotFound:
		return "STATIONID_NOT_FOUND"
	case ErrServiceOffline:
		return "SERVICE_OFFLINE"
	case ErrAccountExpired:
		return "ACCOUNT_EXPIRED"
	case ErrInvalidHash:
		return "INVALID_HASH"
	case ErrInvalidUser:
		return "INVALID_USER"
	case ErrAccountLockout:
		return "ACCOUNT_LOCKOUT"
	case ErrAccountDisabled:
		return "ACCOUNT_DISABLED"
	case ErrTokenExpired:
		return "TOKEN_EXPIRED"
	case ErrMaxLineupChangesReached:
		return "MAX_LINEUP_CHANGES_REACHED"
	case ErrMaxLineups:
		return "MAX_LINEUPS"
	case ErrNoLineups:
		return "NO_LINEUPS"
	case ErrImageNotFound:
		return "IMAGE_NOT_FOUND"
	case ErrInvalidProgramID:
		return "INVALID_PROGRAMID"
	case ErrProgramidQueued:
		return "PROGRAMID_QUEUED"
	case ErrScheduleNotFound:
		return "SCHEDULE_NOT_FOUND"
	case ErrInvalidScheduleRequest:
		return "INVALID_SCHEDULE_REQUEST"
	case ErrScheduleRangeExceeded:
		return "SCHEDULE_RANGE_EXCEEDED"
	case ErrScheduleNotInLineup:
		return "SCHEDULE_NOT_IN_LINEUP"
	case ErrScheduleQueued:
		return "SCHEDULE_QUEUED"
	case ErrHCF:
		return "HCF"
	default:
		return "Unknown ErrorCode(" + strconv.FormatInt(int64(c), 10) + ")"
	}
}

func (c ErrorCode) String() string {
	switch c {
	case ErrOK:
		return "OK"
	case ErrInvalidJSON:
		return "Unable to decode JSON"
	case ErrDeflateRequired:
		return "Did not receive Accept-Encoding: deflate in request."
	case ErrTokenMissing:
		return "Token required but not provided in request header."
	case ErrUnsupportedCommand:
		return "Unsupported command"
	case ErrRequiredActionMissing:
		return "Request is missing an action to take."
	case ErrRequiredRequestMissing:
		return "Did not receive request."
	case ErrRequiredParameterMissingCountry:
		return "In order to search for lineups, you must supply a 3-letter country parameter."
	case ErrRequiredParameterMissingPostalCode:
		return "In order to search for lineups, you must supply a postal code parameter."
	case ErrRequiredParameterMissingMessageID:
		return "In order to delete a message you must supply the messageID."
	case ErrInvalidParameterCountry:
		return "The COUNTRY parameter must be ISO-3166-1 alpha 3. See http://en.wikipedia.org/wiki/ISO_3166-1_alpha-3"
	case ErrInvalidParameterPostalCode:
		return "The POSTALCODE parameter must be valid for the country you are searching. Post message to http://forums.schedulesdirect.org/viewforum.php?f=6 if you are having issues."
	case ErrInvalidParameterFetchType:
		return "You didn't provide a fetchtype I know how to handle."
	case ErrDuplicateLineup:
		return "Lineup already in account."
	case ErrLineupNotFound:
		return "Lineup not in account. Add lineup to account before requesting mapping."
	case ErrUnknownLineup:
		return "Invalid lineup requested. Check your COUNTRY / POSTALCODE combination for validity."
	case ErrInvalidLineupDelete:
		return "Delete of lineup not in account."
	case ErrLineupWrongFormat:
		return "Lineup must be formatted COUNTRY-LINEUP-DEVICE or COUNTRY-OTA-POSTALCODE"
	case ErrInvalidLineup:
		return "The lineup you submitted doesn't exist."
	case ErrLineupDeleted:
		return "The lineup you requested has been deleted from the server."
	case ErrLineupQueued:
		return "The lineup is being generated on the server. Please retry."
	case ErrInvalidCountry:
		return "The country you requested is either mis-typed or does not have valid data."
	case ErrStationIDNotFound:
		return "The stationID you requested is not in any of your lineups."
	case ErrServiceOffline:
		return "Server offline for maintenance."
	case ErrAccountExpired:
		return "Account expired."
	case ErrInvalidHash:
		return "Password hash must be lowercase 40 character sha1_hex of password."
	case ErrInvalidUser:
		return "Invalid username or password."
	case ErrAccountLockout:
		return "Too many login failures. Locked for 15 minutes."
	case ErrAccountDisabled:
		return "Account has been disabled. Please contact Schedules Direct support: admin@schedulesdirect.org for more information."
	case ErrTokenExpired:
		return "Token has expired. Request new token."
	case ErrMaxLineupChangesReached:
		return "Exceeded maximum number of lineup changes for today."
	case ErrMaxLineups:
		return "Exceeded number of lineups for this account."
	case ErrNoLineups:
		return "No lineups have been added to this account."
	case ErrImageNotFound:
		return "Could not find requested image. Post message to http://forums.schedulesdirect.org/viewforum.php?f=6 if you are having issues."
	case ErrInvalidProgramID:
		return "Could not find requested programID. Permanent failure."
	case ErrProgramidQueued:
		return "ProgramID should exist at the server, but doesn't. The server will regenerate the JSON for the program, so your application should retry."
	case ErrScheduleNotFound:
		return "The schedule you requested should be available. Post message to http://forums.schedulesdirect.org/viewforum.php?f=6"
	case ErrInvalidScheduleRequest:
		return "The server can't determine whether your schedule is valid or not. Open a support ticket."
	case ErrScheduleRangeExceeded:
		return "The date that you've requested is outside of the range of the data for that stationID."
	case ErrScheduleNotInLineup:
		return "You have requested a schedule which is not in any of your configured lineups."
	case ErrScheduleQueued:
		return "The schedule you requested has been queued for generation but is not yet ready for download. Retry."
	case ErrHCF:
		return "Unknown error. Open support ticket."
	default:
		return "Unknown ErrorCode(" + strconv.FormatInt(int64(c), 10) + ")"
	}
}
