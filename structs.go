package schedulesdirect

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

// Date is a Schedules Direct specific date format (YYYY-MM-DD) with (Un)MarshalJSON functions.
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
func (p *Date) UnmarshalJSON(text []byte) (err error) {
	strDate := string(text[1:11])

	dateFormat := "2006-01-02"
	if len(strDate) == 4 { // Year only
		dateFormat = "2006"
	}

	v, e := time.Parse(dateFormat, strDate)
	if e != nil {
		return fmt.Errorf("schedulesdirect.Date should be a time, error value is: %s", strDate)
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

// BaseResponse contains the fields that every request is expected to return.
type BaseResponse struct {
	Response string     `json:"response,omitempty"`
	Code     ErrorCode  `json:"code,omitempty"`
	ServerID string     `json:"serverID,omitempty"`
	Message  string     `json:"message,omitempty"`
	DateTime *time.Time `json:"datetime,omitempty"`
}

// Error returns a error string.
func (e *BaseResponse) Error() string {
	msg := e.Message
	if msg == "" {
		msg = e.Response
	}
	return fmt.Sprintf("%s (%d)", msg, e.Code)
}

// A TokenResponse stores the response for token request.
type TokenResponse struct {
	BaseResponse BaseResponse

	Token string `json:"token,omitempty"`
}

// A VersionResponse stores the response for a version request.
type VersionResponse struct {
	BaseResponse BaseResponse

	Client  string `json:"client,omitempty"`
	Version string `json:"version,omitempty"`
}

// A ChangeLineupResponse stores the message after attempting
// to add or delete a lineup.
type ChangeLineupResponse struct {
	BaseResponse BaseResponse

	ChangesRemaining jsonInt `json:"changesRemaining,omitempty"`
}

// A LineupResponse stores the message after requesting
// to list subscribed lineups.
type LineupResponse struct {
	BaseResponse BaseResponse

	Lineups []Lineup `json:"lineups,omitempty"`
}

// A StatusResponse stores the message after requesting system
// status.  SystemStatus[0].Status should be "Online" before proceeding.
type StatusResponse struct {
	Account        *AccountInfo `json:"account,omitempty"`
	Lineups        []Lineup     `json:"lineups,omitempty"`
	LastDataUpdate time.Time    `json:"lastDataUpdate,omitempty"`
	Notifications  []string     `json:"notifications,omitempty"`
	SystemStatus   []Status     `json:"systemStatus,omitempty"`
}

// A StatusError struct stores the error response to a status request.
type StatusError struct {
	BaseResponse BaseResponse

	Token string `json:"token,omitempty"`
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

// A Headend stores the message information for a headend.
type Headend struct {
	Headend   string   `json:"headend,omitempty"`
	Transport string   `json:"transport,omitempty"`
	Location  string   `json:"location,omitempty"`
	Lineups   []Lineup `json:"lineups,omitempty"`
}

// A Lineup stores the message lineup information.
type Lineup struct {
	Lineup    string     `json:"lineup,omitempty"`
	Name      string     `json:"name,omitempty"`
	ID        string     `json:"ID,omitempty"`
	Modified  *time.Time `json:"modified,omitempty"`
	URI       string     `json:"uri,omitempty"`
	IsDeleted bool       `json:"isDeleted,omitempty"`
}

// A BroadcasterInfo stores the information about a broadcaster.
type BroadcasterInfo struct {
	City       string `json:"city,omitempty"`
	State      string `json:"state,omitempty"`
	Postalcode string `json:"postalcode,omitempty"`
	Country    string `json:"country,omitempty"`
}

// A ChannelResponse stores the channel response for a lineup
type ChannelResponse struct {
	Map      []ChannelMap         `json:"map,omitempty"`
	Stations []Station            `json:"stations,omitempty"`
	Metadata *ChannelResponseMeta `json:"metadata,omitempty"`
}

// A ChannelResponseMeta stores the metadata field associated with a channel response
type ChannelResponseMeta struct {
	Lineup     string     `json:"lineup,omitempty"`
	Modified   *time.Time `json:"modified,omitempty"`
	Transport  string     `json:"transport,omitempty"`
	Modulation string     `json:"modulation,omitempty"`
}

// A Station stores the that a station.
type Station struct {
	Affiliate           string           `json:"affiliate,omitempty"`
	Broadcaster         *BroadcasterInfo `json:"broadcaster,omitempty"`
	BroadcastLanguage   []string         `json:"broadcastLanguage,omitempty"`
	CallSign            string           `json:"callsign,omitempty"`
	DescriptionLanguage []string         `json:"descriptionLanguage,omitempty"`
	IsCommercialFree    bool             `json:"isCommercialFree,omitempty"`
	Logo                *StationLogo     `json:"logo,omitempty"`
	Logos               []StationLogo    `json:"stationLogo,omitempty"`
	Name                string           `json:"name,omitempty"`
	StationID           string           `json:"stationID,omitempty"`
	IsRadioStation      bool             `json:"isRadioStation,omitempty"`
}

// A StationLogo stores the information to locate a station logo
type StationLogo struct {
	URL    string `json:"URL,omitempty"`
	Height int    `json:"height,omitempty"`
	Width  int    `json:"width,omitempty"`
	MD5    string `json:"md5,omitempty"`
	Source string `json:"source,omitempty"`
}

// A ChannelMap stores the station id to channel mapping
type ChannelMap struct {
	Channel              string `json:"channel,omitempty"`
	ChannelMajor         int    `json:"channelMajor,omitempty"`
	ChannelMinor         int    `json:"channelMinor,omitempty"`
	DeliverySystem       string `json:"deliverySystem,omitempty"`
	FED                  string `json:"fec,omitempty"`
	FrequencyHertz       int    `json:"frequencyHz,omitempty"`
	LogicalChannelNumber string `json:"logicalChannelNumber,omitempty"`
	MatchType            string `json:"matchType,omitempty"`
	ModulationSystem     string `json:"modulationSystem,omitempty"`
	NetworkID            int    `json:"networkID,omitempty"`
	Polarization         string `json:"polarization,omitempty"`
	ProviderCallSign     string `json:"providerCallsign,omitempty"`
	ServiceID            int    `json:"serviceID,omitempty"`
	StationID            string `json:"stationID,omitempty"`
	SymbolRate           int    `json:"symbolrate,omitempty"`
	TransportID          int    `json:"transportID,omitempty"`
	VirtualChannel       string `json:"virtualChannel,omitempty"`
}

// A Schedule stores the program information for a given stationID
type Schedule struct {
	StationID string        `json:"stationID,omitempty"`
	Metadata  *ScheduleMeta `json:"metadata,omitempty"`
	Programs  []Program     `json:"programs,omitempty"`
}

// A ScheduleMeta stores the metadata information for a schedule
type ScheduleMeta struct {
	Modified  *time.Time `json:"modified,omitempty"`
	MD5       string     `json:"md5,omitempty"`
	StartDate *Date      `json:"startDate,omitempty"`
	EndDate   *Date      `json:"endDate,omitempty"`
	Days      int        `json:"days,omitempty"`
}

// StationScheduleRequest is the payload used to get schedule information for a station as well as last modified information.
type StationScheduleRequest struct {
	StationID string   `json:"stationID,omitempty"`
	Dates     []string `json:"date,omitempty"`
}

// LastModifiedEntry contains information about the last modification of a station schedule.
type LastModifiedEntry struct {
	LastModified *time.Time `json:"lastModified,omitempty"`
	MD5          string     `json:"md5,omitempty"`
}

// LanguageCrossReference provides translated titles and descriptions for a program.
type LanguageCrossReference struct {
	BaseResponse BaseResponse

	DescriptionLanguage     string `json:"descriptionLanguage,omitempty"`
	DescriptionLanguageName string `json:"descriptionLanguageName,omitempty"`
	MD5                     string `json:"md5,omitempty"`
	ProgramID               string `json:"programID,omitempty"`
	TitleLanguage           string `json:"titleLanguage,omitempty"`
	TitleLanguageName       string `json:"titleLanguageName,omitempty"`
}

// A StillRunningResponse describes the current real time state of a program.
type StillRunningResponse struct {
	BaseResponse BaseResponse

	EventStartDateTime *time.Time `json:"eventStartDateTime,omitempty"`
	IsComplete         bool       `json:"isComplete,omitempty"`
	ProgramID          string     `json:"programID,omitempty"`
	Result             struct {
		AwayTeam *Team `json:"awayTeam,omitempty"`
		HomeTeam *Team `json:"homeTeam,omitempty"`
	} `json:"result,omitempty"`
}
