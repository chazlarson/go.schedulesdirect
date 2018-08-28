package schedulesdirect

import (
	"fmt"
	"time"
)

// Date is a Schedules Direct specific date format (YYYY-MM-DD) with (Un)MarshalJSON functions.
type Date struct {
	time.Time
	fmt string
}

// MarshalJSON formats the underlying time.Time to Schedule Direct's format.
func (p Date) MarshalJSON() ([]byte, error) {
	t := p.Time
	str := "\"" + t.Format(p.fmt) + "\""

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
	*p = Date{v, dateFormat}
	return nil
}

// BaseResponse contains the fields that every request is expected to return.
type BaseResponse struct {
	Response string    `json:"response"`
	Code     ErrorCode `json:"code"`
	ServerID string    `json:"serverID"`
	Message  string    `json:"message"`
	DateTime time.Time `json:"datetime"`
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

	Token string `json:"token"`
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

	ChangesRemaining int `json:"changesRemaining"`
}

// A LineupResponse stores the message after requesting
// to list subscribed lineups.
type LineupResponse struct {
	BaseResponse BaseResponse

	Lineups []Lineup `json:"lineups"`
}

// A StatusResponse stores the message after requesting system
// status.  SystemStatus[0].Status should be "Online" before proceeding.
type StatusResponse struct {
	Account        AccountInfo `json:"account"`
	Lineups        []Lineup    `json:"lineups"`
	LastDataUpdate time.Time   `json:"lastDataUpdate"`
	Notifications  []string    `json:"notifications"`
	SystemStatus   []Status    `json:"systemStatus"`
}

// A StatusError struct stores the error response to a status request.
type StatusError struct {
	BaseResponse BaseResponse

	Token string `json:"token"`
}

// A Status stores the message system status information
// usually as part of a StatusResponse.
type Status struct {
	Date    time.Time `json:"date"`
	Status  string    `json:"status"`
	Details string    `json:"details"`
}

// An AccountInfo stores the message account information
// usually as part of a StatusResponse.
type AccountInfo struct {
	Expires    string   `json:"expires"`
	Messages   []string `json:"messages"`
	MaxLineups int      `json:"maxLineups"`
}

// A Headend stores the message information for a headend.
type Headend struct {
	Headend   string   `json:"headend"`
	Transport string   `json:"transport"`
	Location  string   `json:"location"`
	Lineups   []Lineup `json:"lineups"`
}

// A Lineup stores the message lineup information.
type Lineup struct {
	Lineup    string    `json:"lineup,omitempty"`
	Name      string    `json:"name,omitempty"`
	ID        string    `json:"ID,omitempty"`
	Modified  time.Time `json:"modified,omitempty"`
	URI       string    `json:"uri"`
	IsDeleted bool      `json:"isDeleted,omitempty"`
}

// A BroadcasterInfo stores the information about a broadcaster.
type BroadcasterInfo struct {
	City       string `json:"city"`
	State      string `json:"state"`
	Postalcode string `json:"postalcode"`
	Country    string `json:"country"`
}

// A ChannelResponse stores the channel response for a lineup
type ChannelResponse struct {
	Map      []ChannelMap        `json:"map"`
	Stations []Station           `json:"stations"`
	Metadata ChannelResponseMeta `json:"metadata"`
}

// A ChannelResponseMeta stores the metadata field associated with a channel response
type ChannelResponseMeta struct {
	Lineup     string    `json:"lineup"`
	Modified   time.Time `json:"modified"`
	Transport  string    `json:"transport"`
	Modulation string    `json:"modulation"`
}

// A Station stores the that a station.
type Station struct {
	Affiliate           string          `json:"affiliate"`
	Broadcaster         BroadcasterInfo `json:"broadcaster"`
	BroadcastLanguage   []string        `json:"broadcastLanguage"`
	CallSign            string          `json:"callsign"`
	DescriptionLanguage []string        `json:"descriptionLanguage"`
	IsCommercialFree    bool            `json:"isCommercialFree"`
	Logo                StationLogo     `json:"logo"`
	Logos               []StationLogo   `json:"stationLogo"`
	Name                string          `json:"name"`
	StationID           string          `json:"stationID"`
	IsRadioStation      bool            `json:"isRadioStation"`
}

// A StationLogo stores the information to locate a station logo
type StationLogo struct {
	URL    string `json:"URL"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
	MD5    string `json:"md5"`
	Source string `json:"source"`
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
	StationID string       `json:"stationID"`
	Metadata  ScheduleMeta `json:"metadata"`
	Programs  []Program    `json:"programs"`
}

// A ScheduleMeta stores the metadata information for a schedule
type ScheduleMeta struct {
	Modified  time.Time `json:"modified"`
	MD5       string    `json:"md5"`
	StartDate Date      `json:"startDate,omitempty"`
	EndDate   Date      `json:"endDate,omitempty"`
	Days      int       `json:"days,omitempty"`
}

// StationScheduleRequest is the payload used to get schedule information for a station as well as last modified information.
type StationScheduleRequest struct {
	StationID string   `json:"stationID"`
	Dates     []string `json:"date,omitempty"`
}

// LastModifiedEntry contains information about the last modification of a station schedule.
type LastModifiedEntry struct {
	LastModified time.Time `json:"lastModified"`
	MD5          string    `json:"md5"`
}

// LanguageCrossReference provides translated titles and descriptions for a program.
type LanguageCrossReference struct {
	BaseResponse BaseResponse

	DescriptionLanguage     string `json:"descriptionLanguage"`
	DescriptionLanguageName string `json:"descriptionLanguageName"`
	MD5                     string `json:"md5"`
	ProgramID               string `json:"programID"`
	TitleLanguage           string `json:"titleLanguage"`
	TitleLanguageName       string `json:"titleLanguageName"`
}

// A StillRunningResponse describes the current real time state of a program.
type StillRunningResponse struct {
	BaseResponse BaseResponse

	EventStartDateTime time.Time `json:"eventStartDateTime"`
	IsComplete         bool      `json:"isComplete"`
	ProgramID          string    `json:"programID"`
	Result             struct {
		AwayTeam Team `json:"awayTeam"`
		HomeTeam Team `json:"homeTeam"`
	} `json:"result"`
}
