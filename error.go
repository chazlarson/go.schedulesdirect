package schedulesdirect

import "strconv"

// A ErrorCode is an unsigned 32-bit error code as defined in the Schedules Direct JSON specification.
//
// This implementation was inspired by gRPC's codes package: https://github.com/grpc/grpc-go/tree/master/codes.
type ErrorCode uint32

const (

	// OK is returned when OK
	OK ErrorCode = 0

	// InvalidJSON is returned when unable to decode JSON.
	InvalidJSON ErrorCode = 1001

	// DeflateRequired is returned when did not receive Accept-Encoding: deflate in request.
	DeflateRequired ErrorCode = 1002

	// TokenMissing is returned when token required but not provided in request header.
	TokenMissing ErrorCode = 1004

	// UnsupportedCommand is returned when unsupported command.
	UnsupportedCommand ErrorCode = 2000

	// RequiredActionMissing is returned when request is missing an action to take.
	RequiredActionMissing ErrorCode = 2001

	// RequiredRequestMissing is returned when did not receive request.
	RequiredRequestMissing ErrorCode = 2002

	// RequiredParameterMissingCountry is returned when in order to search for lineups, you must supply a 3-letter country parameter.
	RequiredParameterMissingCountry ErrorCode = 2004

	// RequiredParameterMissingPostalCode is returned when in order to search for lineups, you must supply a postal code parameter.
	RequiredParameterMissingPostalCode ErrorCode = 2005

	// RequiredParameterMissingMessageID is returned when in order to delete a message you must supply the messageID.
	RequiredParameterMissingMessageID ErrorCode = 2006

	// InvalidParameterCountry is returned when the COUNTRY parameter was not a ISO-3166-1 alpha 3 code.
	// See http://en.wikipedia.org/wiki/ISO_3166-1_alpha-3 for more information.
	InvalidParameterCountry ErrorCode = 2050

	// InvalidParameterPostalCode is returned when the POSTALCODE parameter must was invalid for the country searched.
	// Post message to http://forums.schedulesdirect.org/viewforum.php?f=6 if you are having issues.
	InvalidParameterPostalCode ErrorCode = 2051

	// InvalidParameterFetchType is returned when you didn't provide a fetchtype Schedules Direct knows how to handle.
	InvalidParameterFetchType ErrorCode = 2052

	// DuplicateLineup is returned when a lineup already existed in account.
	DuplicateLineup ErrorCode = 2100

	// LineupNotFound is returned when lineup is not in an account.
	// Add lineup to account before requesting mapping.
	LineupNotFound ErrorCode = 2101

	// UnknownLineup is returned when invalid lineup requested.
	// Check your COUNTRY / POSTALCODE combination for validity.
	UnknownLineup ErrorCode = 2102

	// InvalidLineupDelete is returned when delete of lineup not in account.
	InvalidLineupDelete ErrorCode = 2103

	// LineupWrongFormat is returned when lineup must be formatted COUNTRY-LINEUP-DEVICE or COUNTRY-OTA-POSTALCODE.
	LineupWrongFormat ErrorCode = 2104

	// InvalidLineup is returned when the lineup you submitted doesn't exist.
	InvalidLineup ErrorCode = 2105

	// LineupDeleted is returned when the lineup you requested has been deleted from the server.
	LineupDeleted ErrorCode = 2106

	// LineupQueued is returned when the lineup is being generated on the server.
	// Please retry.
	LineupQueued ErrorCode = 2107

	// InvalidCountry is returned when the country you requested is either mis-typed or does not have valid data.
	InvalidCountry ErrorCode = 2108

	// StationIDNotFound is returned when the stationID you requested is not in any of your lineups.
	StationIDNotFound ErrorCode = 2200

	// ServiceOffline is returned when server offline for maintenance.
	ServiceOffline ErrorCode = 3000

	// AccountExpired is returned when account expired.
	AccountExpired ErrorCode = 4001

	// InvalidHash is returned when password hash must be lowercase 40 character sha1_hex of password.
	InvalidHash ErrorCode = 4002

	// InvalidUser is returned when invalid username or password.
	InvalidUser ErrorCode = 4003

	// AccountLockout is returned when too many login failures.
	// Locked for 15 minutes.
	AccountLockout ErrorCode = 4004

	// AccountDisabled is returned when account has been disabled.
	// Please contact Schedules Direct support: admin@schedulesdirect.org for more information.
	AccountDisabled ErrorCode = 4005

	// TokenExpired is returned when token has expired.
	// Request new token.
	TokenExpired ErrorCode = 4006

	// MaxLineupChangesReached is returned when exceeded maximum number of lineup changes for today.
	MaxLineupChangesReached ErrorCode = 4100

	// MaxLineups is returned when exceeded number of lineups for this account.
	MaxLineups ErrorCode = 4101

	// NoLineups is returned when no lineups have been added to this account.
	NoLineups ErrorCode = 4102

	// ImageNotFound is returned when could not find requested image.
	// Post message to http://forums.schedulesdirect.org/viewforum.php?f=6 if you are having issues.
	ImageNotFound ErrorCode = 5000

	// InvalidProgramID is returned when could not find requested programID.
	// Permanent failure.
	InvalidProgramID ErrorCode = 6000

	// ProgramidQueued is returned when programID should exist at the server, but doesn't.
	// The server will regenerate the JSON for the program, so your application should retry.
	ProgramidQueued ErrorCode = 6001

	// ScheduleNotFound is returned when the schedule you requested should be available.
	// Post message to http://forums.schedulesdirect.org/viewforum.php?f=6
	ScheduleNotFound ErrorCode = 7000

	// InvalidScheduleRequest is returned when the server can't determine whether your schedule is valid or not.
	// Open a support ticket.
	InvalidScheduleRequest ErrorCode = 7010

	// ScheduleRangeExceeded is returned when the date that you've requested is outside of the range of the data for that stationID.
	ScheduleRangeExceeded ErrorCode = 7020

	// ScheduleNotInLineup is returned when you have requested a schedule which is not in any of your configured lineups.
	ScheduleNotInLineup ErrorCode = 7030

	// ScheduleQueued is returned when the schedule you requested has been queued for generation but is not yet ready for download.
	// Retry.
	ScheduleQueued ErrorCode = 7100

	// HCF is returned when Unknown error.
	// Open support ticket.
	HCF ErrorCode = 9999
)

func (c ErrorCode) String() string {
	switch c {
	case OK:
		return "OK"
	case InvalidJSON:
		return "Unable to decode JSON"
	case DeflateRequired:
		return "Did not receive Accept-Encoding: deflate in request."
	case TokenMissing:
		return "Token required but not provided in request header."
	case UnsupportedCommand:
		return "Unsupported command"
	case RequiredActionMissing:
		return "Request is missing an action to take."
	case RequiredRequestMissing:
		return "Did not receive request."
	case RequiredParameterMissingCountry:
		return "In order to search for lineups, you must supply a 3-letter country parameter."
	case RequiredParameterMissingPostalCode:
		return "In order to search for lineups, you must supply a postal code parameter."
	case RequiredParameterMissingMessageID:
		return "In order to delete a message you must supply the messageID."
	case InvalidParameterCountry:
		return "The COUNTRY parameter must be ISO-3166-1 alpha 3. See http://en.wikipedia.org/wiki/ISO_3166-1_alpha-3"
	case InvalidParameterPostalCode:
		return "The POSTALCODE parameter must be valid for the country you are searching. Post message to http://forums.schedulesdirect.org/viewforum.php?f=6 if you are having issues."
	case InvalidParameterFetchType:
		return "You didn't provide a fetchtype I know how to handle."
	case DuplicateLineup:
		return "Lineup already in account."
	case LineupNotFound:
		return "Lineup not in account. Add lineup to account before requesting mapping."
	case UnknownLineup:
		return "Invalid lineup requested. Check your COUNTRY / POSTALCODE combination for validity."
	case InvalidLineupDelete:
		return "Delete of lineup not in account."
	case LineupWrongFormat:
		return "Lineup must be formatted COUNTRY-LINEUP-DEVICE or COUNTRY-OTA-POSTALCODE"
	case InvalidLineup:
		return "The lineup you submitted doesn't exist."
	case LineupDeleted:
		return "The lineup you requested has been deleted from the server."
	case LineupQueued:
		return "The lineup is being generated on the server. Please retry."
	case InvalidCountry:
		return "The country you requested is either mis-typed or does not have valid data."
	case StationIDNotFound:
		return "The stationID you requested is not in any of your lineups."
	case ServiceOffline:
		return "Server offline for maintenance."
	case AccountExpired:
		return "Account expired."
	case InvalidHash:
		return "Password hash must be lowercase 40 character sha1_hex of password."
	case InvalidUser:
		return "Invalid username or password."
	case AccountLockout:
		return "Too many login failures. Locked for 15 minutes."
	case AccountDisabled:
		return "Account has been disabled. Please contact Schedules Direct support: admin@schedulesdirect.org for more information."
	case TokenExpired:
		return "Token has expired. Request new token."
	case MaxLineupChangesReached:
		return "Exceeded maximum number of lineup changes for today."
	case MaxLineups:
		return "Exceeded number of lineups for this account."
	case NoLineups:
		return "No lineups have been added to this account."
	case ImageNotFound:
		return "Could not find requested image. Post message to http://forums.schedulesdirect.org/viewforum.php?f=6 if you are having issues."
	case InvalidProgramID:
		return "Could not find requested programID. Permanent failure."
	case ProgramidQueued:
		return "ProgramID should exist at the server, but doesn't. The server will regenerate the JSON for the program, so your application should retry."
	case ScheduleNotFound:
		return "The schedule you requested should be available. Post message to http://forums.schedulesdirect.org/viewforum.php?f=6"
	case InvalidScheduleRequest:
		return "The server can't determine whether your schedule is valid or not. Open a support ticket."
	case ScheduleRangeExceeded:
		return "The date that you've requested is outside of the range of the data for that stationID."
	case ScheduleNotInLineup:
		return "You have requested a schedule which is not in any of your configured lineups."
	case ScheduleQueued:
		return "The schedule you requested has been queued for generation but is not yet ready for download. Retry."
	case HCF:
		return "Unknown error. Open support ticket."
	default:
		return "Unknown ErrorCode(" + strconv.FormatInt(int64(c), 10) + ")"
	}
}
