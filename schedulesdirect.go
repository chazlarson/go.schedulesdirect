// Package schedulesdirect provides structs and functions to interact with
// the SchedulesDirect JSON listings service in Go.
package schedulesdirect

import (
	"bytes"
	"compress/gzip"
	"crypto/sha1" // #nosec
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

//The proper SchedulesDirect JSON Service workflow is as follows...
//------Once the client is in a steady state:-------
//-Obtain a token
//-Obtain the current status.
//-If the system status is "OFFLINE" then disconnect; all further processing
//   will be rejected at the server. A client should not attempt to reconnect
//   for 1 hour.
//-Check the status object and determine if any headends on the server have
//   newer "modified" dates than the one that is on the client. If yes, download
//   the updated lineup for that headend.
//-If there are no changes to the headends, send a request to the server for
//   the MD5 hashes of the schedules that you are interested in. If the MD5
//   hash for the schedule is the same as you have locally cached from your
//   last download, then the schedule on the server hasn't changed and your
//   client should disconnect.
//-If the MD5 hash for the schedule is different, then download the schedules
//   that have different hashes.
//-Honor the nextScheduled time in the status object; if your client connects
//   during server-side data processing, the nextScheduled time will be
//   "closer", however reconnecting while server-side data is being processed
//   will not result in newer data.
//-Parse the schedule, determine if the MD5 of the program for a particular
//   timeslot has changed. If the program ID for a timeslot is the same, but
//   the MD5 has changed, this means that some sort of metadata for that
//   program has been updated.
//-Request the "delta" program id's as determined through the MD5 values.

// Some constants for use in the library
var (
	APIVersion     = "20141201"
	DefaultBaseURL = "https://json.schedulesdirect.org/"
	UserAgent      = "go.schedulesdirect (Go-http-client/1.1)"
)

// BaseResponse contains the fields that every request is expected to return.
type BaseResponse struct {
	Response string    `json:"response"`
	Code     int       `json:"code"`
	ServerID string    `json:"serverID"`
	Message  string    `json:"message"`
	DateTime time.Time `json:"datetime"`
}

// Error returns a error string.
func (e *BaseResponse) Error() string {
	return e.Message
}

// A TokenResponse stores the SD json response message for token request.
type TokenResponse struct {
	BaseResponse

	Token string `json:"token"`
}

// A VersionResponse stores the SD json response message for a version request.
type VersionResponse struct {
	BaseResponse

	Version string `json:"version,omitempty"`
}

// A ChangeLineupResponse stores the SD json message returned after attempting
// to add or delete a lineup.
type ChangeLineupResponse struct {
	BaseResponse

	ChangesRemaining int `json:"changesRemaining"`
}

// A LineupResponse stores the SD json message returned after requesting
// to list subscribed lineups.
type LineupResponse struct {
	BaseResponse

	Lineups []Lineup `json:"lineups"`
}

// A StatusResponse stores the SD json message returned after requesting system
// status.  SystemStatus[0].Status should be "Online" before proceeding.
type StatusResponse struct {
	BaseResponse

	Account        AccountInfo `json:"account"`
	Lineups        []Lineup    `json:"lineups"`
	LastDataUpdate string      `json:"lastDataUpdate"`
	Notifications  []string    `json:"notifications"`
	SystemStatus   []Status    `json:"systemStatus"`
}

// A StatusError struct stores the error response to a status request.
type StatusError struct {
	BaseResponse

	Token string `json:"token"`
}

// A Status stores the SD json message containing system status information
// usually as part of a StatusResponse.
type Status struct {
	Date    string `json:"date"`
	Status  string `json:"status"`
	Details string `json:"details"`
}

// An AccountInfo stores the SD json message containing account information
// usually as part of a StatusResponse.
type AccountInfo struct {
	Expires                  string   `json:"expires"`
	Messages                 []string `json:"messages"`
	MaxLineups               int      `json:"maxLineups"`
	NextSuggestedConnectTime string   `json:"nextSuggestedConnectTime"`
}

// A Headend stores the SD json message containing information for a headend.
type Headend struct {
	Headend   string   `json:"headend"`
	Transport string   `json:"transport"`
	Location  string   `json:"location"`
	Lineups   []Lineup `json:"lineups"`
}

// A Lineup stores the SD json message containing lineup information.
type Lineup struct {
	Lineup    string `json:"lineup,omitempty"`
	Name      string `json:"name,omitempty"`
	ID        string `json:"ID,omitempty"`
	Modified  string `json:"modified,omitempty"`
	URI       string `json:"uri"`
	IsDeleted bool   `json:"isDeleted,omitempty"`
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

// A Station stores the SD json that describes a station.
type Station struct {
	Callsign            string          `json:"callsign"`
	Affiliate           string          `json:"affiliate"`
	IsCommercialFree    bool            `json:"isCommercialFree"`
	Name                string          `json:"name"`
	Broadcaster         BroadcasterInfo `json:"broadcaster"`
	BroadcastLanguage   []string        `json:"broadcastLanguage"`
	DescriptionLanguage []string        `json:"descriptionLanguage "`
	Logo                StationLogo     `json:"logo"`
	StationID           string          `json:"stationID"`
}

// A StationLogo stores the information to locate a station logo
type StationLogo struct {
	URL    string `json:"URL"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
	MD5    string `json:"md5"`
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
	ProviderCallsign     string `json:"providerCallsign,omitempty"`
	ServiceID            int    `json:"serviceID,omitempty"`
	StationID            string `json:"stationID,omitempty"`
	Symbolrate           int    `json:"symbolrate,omitempty"`
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
	Modified  string `json:"modified"`
	MD5       string `json:"md5"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Days      int    `json:"days"`
}

// A SyndicationType stores syndication information for a program
type SyndicationType struct {
	Source string `json:"source"`
	Type   string `json:"type"`
}

// A Program stores the information to describing a single television program.
type Program struct {
	ProgramID           string          `json:"programID,omitempty"`
	AirDateTime         time.Time       `json:"airDateTime,omitempty"`
	MD5                 string          `json:"md5,omitempty"`
	Duration            int             `json:"duration,omitempty"`
	New                 bool            `json:"new,omitempty"`
	CableInTheClassroom bool            `json:"cableInTheClassRoom,omitempty"`
	Catchup             bool            `json:"catchup,omitempty"`   // - typically only found outside of North America
	Continued           bool            `json:"continued,omitempty"` // - typically only found outside of North America
	Education           bool            `json:"educational,omitempty"`
	JoinedInProgress    bool            `json:"joinedInProgress,omitempty"`
	LeftInProgress      bool            `json:"leftInProgress,omitempty"`
	Premiere            bool            `json:"premiere,omitempty"`          //- Should only be found in Miniseries and Movie program types.
	ProgramBreak        bool            `json:"programBreak,omitempty"`      // - Program stops and will restart later (frequently followed by a continued). Typically only found outside of North America.
	Repeat              bool            `json:"repeat,omitempty"`            // - An encore presentation. Repeat should only be found on a second telecast of sporting events.
	Signed              bool            `json:"signed,omitempty"`            //- Program has an on-screen person providing sign-language translation.
	SubjectToBlackout   bool            `json:"subjectToBlackout,omitempty"` //subjectToBlackout
	TimeApproximate     bool            `json:"timeApproximate,omitempty"`
	AudioProperties     []string        `json:"audioProperties,omitempty"`
	Syndication         SyndicationType `json:"syndication,omitempty"`
	Ratings             []Rating        `json:"ratings,omitempty"`
	ProgramPart         Part            `json:"multipart,omitempty"`
	VideoProperties     []string        `json:"videoProperties,omitempty"`
}

// A Rating stores ratings board information for a program
type Rating struct {
	Body string `json:"body"`
	Code string `json:"code"`
}

// Title contains the title of a program.
type Title struct {
	Title120 string `json:"title120,omitempty"`
}

// EventDetails indicates the type of program.
type EventDetails struct {
	SubType *string `json:"subType,omitempty"`
}

// Metadata stores meta information for a program.
type Metadata struct {
	Episode       int `json:"episode,omitzero"`
	Season        int `json:"season,omitzero"`
	TotalEpisodes int `json:"totalEpisodes,omitempty"`
	TotalSeasons  int `json:"totalSeasons,omitempty"`
}

// A ProgramInfo type stores program information for a program
type ProgramInfo struct {
	ProgramID       string                   `json:"programID,omitempty"`
	Titles          []Title                  `json:"titles,omitempty"`
	EventDetails    EventDetails             `json:"eventDetails,omitempty"`
	Descriptions    map[string][]Description `json:"descriptions,omitempty"`
	OriginalAirDate string                   `json:"originalAirDate,omitempty"`
	Genres          []string                 `json:"genres,omitempty"`
	EpisodeTitle150 string                   `json:"episodeTitle150,omitempty"`
	Metadata        []map[string]Metadata    `json:"metadata,omitempty"`
	Keywords        map[string][]string      `json:"keyWords,omitempty"`
	Movie           Movie                    `json:"movie,omitempty"`
	Cast            []Person                 `json:"cast,omitempty"`
	Crew            []Person                 `json:"crew,omitempty"`
	ContentRating   []Rating                 `json:"contentRating,omitempty"`
	EntityType      string                   `json:"entityType,omitempty"`
	ShowType        string                   `json:"showType,omitempty"`
	HasImageArtWork bool                     `json:"hasImageArtwork,omitempty"`
	MD5             string                   `json:"md5,omitempty"`
}

// A MovieQualityRating describes ratings for the quality of a movie.
type MovieQualityRating struct {
	Increment   string `json:"increment,omitempty"`
	MaxRating   string `json:"maxRating,omitempty"`
	MinRating   string `json:"minRating,omitempty"`
	Rating      string `json:"rating,omitempty"`
	RatingsBody string `json:"ratingsBody,omitmepty"`
}

// A Movie type stores information about a movie
type Movie struct {
	Duration      int                  `json:"duration,omitempty"`
	Year          string               `json:"year,omitempty"`
	QualityRating []MovieQualityRating `json:"qualityRating,omitempty"`
}

// Person stores information for an acting credit.
type Person struct {
	PersonID      string `json:"personId,omitmepty"`
	NameID        string `json:"nameId,omitempty"`
	Name          string `json:"name,omitempty"`
	Role          string `json:"role,omitempty"`
	CharacterName string `json:"characterName,omitempty"`
	BillingOrder  string `json:"billingOrder,omitempty"`
}

// Description helps store the descriptions for programs
type Description struct {
	DescriptionLanguage string `json:"descriptionLanguage"`
	Description         string `json:"description"`
}

// Part stores the information for a part
type Part struct {
	PartNumber int `json:"partNumber"`
	TotalParts int `json:"totalParts"`
}

// StationScheduleRequest is the payload used to get schedule information for a station as well as last modified information.
type StationScheduleRequest struct {
	StationID string   `json:"stationID"`
	Dates     []string `json:"dates,omitempty"`
}

// LastModifiedEntry contains information about the last modification of a station schedule.
type LastModifiedEntry struct {
	Code         int       `json:"code"`
	LastModified time.Time `json:"lastModified"`
	MD5          string    `json:"md5"`
	Message      string    `json:"message"`
}

// ProgramDescription provides a generic description of a program.
type ProgramDescription struct {
	Code            int    `json:"code"`
	Description100  string `json:"description100"`
	Description1000 string `json:"description1000"`
}

// LanguageCrossReference provides translated titles and descriptions for a program.
type LanguageCrossReference struct {
	BaseResponse

	DescriptionLanguage     string `json:"descriptionLanguage"`
	DescriptionLanguageName string `json:"descriptionLanguageName"`
	MD5                     string `json:"md5"`
	ProgramID               string `json:"programID"`
	TitleLanguage           string `json:"titleLanguage"`
	TitleLanguageName       string `json:"titleLanguageName"`
}

// A StillRunningResponse describes the current real time state of a program.
type StillRunningResponse struct {
	BaseResponse

	EventStartDateTime string `json:"eventStartDateTime"`
	IsComplete         bool   `json:"isComplete"`
	ProgramID          string `json:"programID"`
	Result             struct {
		AwayTeam struct {
			Name  string `json:"name"`
			Score string `json:"score"`
		} `json:"awayTeam"`
		HomeTeam struct {
			Name  string `json:"name"`
			Score string `json:"score"`
		} `json:"homeTeam"`
	} `json:"result"`
}

// ProgramArtwork describes a single piece of artwork related to a program.
type ProgramArtwork struct {
	Aspect   string            `json:"aspect"`
	Category string            `json:"category"`
	Height   string            `json:"height"`
	Primary  string            `json:"primary"`
	Size     string            `json:"size"`
	Text     string            `json:"text"`
	Tier     string            `json:"tier"`
	URI      string            `json:"uri"`
	Width    string            `json:"width"`
	Caption  map[string]string `json:"caption"`
}

// ProgramArtworkResponse is a container struct for artwork relating to a program.
type ProgramArtworkResponse struct {
	ProgramID string           `json:"programID"`
	Artwork   []ProgramArtwork `json:"data"`
}

// Client type
type Client struct {
	// The Base URL for SD requests
	BaseURL *url.URL

	// Our HTTP client to communicate with SD
	HTTP *http.Client

	// The token
	Token string
}

// NewClient returns a new SD API client. Uses http.DefaultClient if no http.Client is set.
func NewClient(username string, password string) (*Client, error) {
	baseURL, parseErr := url.Parse(DefaultBaseURL)
	if parseErr != nil {
		return nil, parseErr
	}
	c := &Client{HTTP: http.DefaultClient, BaseURL: baseURL}
	token, tokenErr := c.GetToken(username, password)
	if tokenErr != nil {
		return nil, tokenErr
	}
	c.Token = token
	return c, nil
}

// encryptPassword returns the sha1 hex encoding of the string argument
func encryptPassword(password string) (string, error) {
	encoded := sha1.New() // #nosec
	if _, writeErr := encoded.Write([]byte(password)); writeErr != nil {
		return "", writeErr
	}
	return hex.EncodeToString(encoded.Sum(nil)), nil
}

// GetToken returns a session token if the supplied username/password
// successfully authenticate.
func (c Client) GetToken(username string, password string) (string, error) {
	url := fmt.Sprint(DefaultBaseURL, APIVersion, "/token")

	// encrypt the password
	sha1hexPW, encryptError := encryptPassword(password)
	if encryptError != nil {
		return "", encryptError
	}

	js, jsErr := json.Marshal(map[string]string{"username": username, "password": sha1hexPW})
	if jsErr != nil {
		return "", jsErr
	}

	// setup the request
	req, httpErr := http.NewRequest("POST", url, bytes.NewBuffer(js))
	if httpErr != nil {
		return "", httpErr
	}

	// perform the POST
	_, data, err := c.sendRequest(req)
	if err != nil {
		return "", err
	}

	// create a TokenResponse struct, return if err
	r := new(TokenResponse)

	// decode the response body into the new TokenResponse struct
	err = json.Unmarshal(data, &r)
	if err != nil {
		return "", err
	}

	// return the token string
	return r.Token, nil
}

// GetStatus returns a StatusResponse for this account.
func (c Client) GetStatus() (*StatusResponse, error) {
	url := fmt.Sprint(DefaultBaseURL, APIVersion, "/status")
	s := new(StatusResponse)

	req, httpErr := http.NewRequest("GET", url, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	_, data, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &s)
	return s, err
}

// AddLineup adds the given lineup uri to the users SchedulesDirect account.
func (c Client) AddLineup(lineupURI string) (*ChangeLineupResponse, error) {
	url := fmt.Sprint(DefaultBaseURL, APIVersion, lineupURI)

	req, httpErr := http.NewRequest("PUT", url, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	_, data, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	r := &ChangeLineupResponse{}

	err = json.Unmarshal(data, &r)
	return r, err
}

// DeleteLineup deletes the given lineup uri from the users SchedulesDirect account.
func (c Client) DeleteLineup(lineupURI string) (*ChangeLineupResponse, error) {
	url := fmt.Sprint(DefaultBaseURL, APIVersion, lineupURI)

	req, httpErr := http.NewRequest("DELETE", url, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	_, data, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	r := &ChangeLineupResponse{}

	err = json.Unmarshal(data, &r)
	return r, err
}

// AutomapLineup accepts the "lineup.json" output as a byte slice from SiliconDust's HDHomerun devices.
// It then runs a comparison against ScheduleDirect's database and returns potential lineup matches.
func (c Client) AutomapLineup(hdhrLineupJSON []byte) (map[string]int, error) {
	url := fmt.Sprint(DefaultBaseURL, APIVersion, "/map/lineup")

	req, httpErr := http.NewRequest("POST", url, bytes.NewBuffer(hdhrLineupJSON))
	if httpErr != nil {
		return nil, httpErr
	}

	_, data, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	matches := make(map[string]int)

	err = json.Unmarshal(data, &matches)
	return matches, err
}

// SubmitLineup should be called if AutomapLineup doesn't return candidates after you identify
// the lineup you were trying to find via automapping.
func (c Client) SubmitLineup(hdhrLineupJSON []byte, lineupID string) (*BaseResponse, error) {
	url := fmt.Sprint(DefaultBaseURL, APIVersion, "/map/lineup/", lineupID)

	req, httpErr := http.NewRequest("POST", url, bytes.NewBuffer(hdhrLineupJSON))
	if httpErr != nil {
		return nil, httpErr
	}

	_, data, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	baseR := &BaseResponse{}

	err = json.Unmarshal(data, &baseR)

	return baseR, err
}

// GetHeadends returns the map of headends for the given country and postal code.
func (c Client) GetHeadends(countryCode, postalCode string) ([]Headend, error) {
	uriPart := fmt.Sprintf("/headends?country=%s&postalcode=%s", countryCode, postalCode)
	url := fmt.Sprint(DefaultBaseURL, APIVersion, uriPart)

	req, httpErr := http.NewRequest("GET", url, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	_, data, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	h := []Headend{}

	err = json.Unmarshal(data, &h)
	if err != nil {
		return nil, err
	}
	return h, nil
}

// GetChannels returns the channels in a given lineup
func (c Client) GetChannels(lineupURI string) (*ChannelResponse, error) {
	url := fmt.Sprint(DefaultBaseURL, APIVersion, lineupURI)

	req, httpErr := http.NewRequest("GET", url, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	_, data, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	h := new(ChannelResponse)

	err = json.Unmarshal(data, &h)
	return h, err
}

// GetSchedules returns the set of schedules requested.  As a whole the response is not valid json but each individual line is valid.
func (c Client) GetSchedules(requests []StationScheduleRequest) ([]Schedule, error) {
	url := fmt.Sprint(DefaultBaseURL, APIVersion, "/schedules")

	js, jsErr := json.Marshal(requests)
	if jsErr != nil {
		return nil, jsErr
	}

	//setup the request
	req, httpErr := http.NewRequest("POST", url, bytes.NewBuffer(js))
	if httpErr != nil {
		return nil, httpErr
	}

	_, data, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	var h []Schedule

	err = json.Unmarshal(data, &h)
	return h, err
}

// GetProgramInfo returns the set of program details for the given set of programs
func (c Client) GetProgramInfo(programIDs []string) ([]ProgramInfo, error) {
	if len(programIDs) > 5000 {
		return nil, errors.New("you may only request at most 5000 program IDs per request, please lower your request amount")
	}
	url := fmt.Sprint(DefaultBaseURL, APIVersion, "/programs")

	js, jsErr := json.Marshal(programIDs)
	if jsErr != nil {
		return nil, jsErr
	}

	// setup the request
	req, httpErr := http.NewRequest("POST", url, bytes.NewBuffer(js))
	if httpErr != nil {
		return nil, httpErr
	}
	req.Header.Set("Accept-Encoding", "deflate,gzip")

	resp, _, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	// Your client must send an Accept-Encoding that has "deflate,gzip" in it, even though the response will be gzip'ed.
	// This is due to an implementation bug in 20140530 which will be fixed in 20141201.
	//
	// Not actually fixed yet and Go disables automatic decompression if Accept-Encoding is set, so we are stuck doing the decompression ourselves.
	var reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		readerG, errG := gzip.NewReader(reader)
		if errG == nil {
			reader = readerG
		} else {
			return nil, errG
		}
	}

	// create the programs slice
	var allPrograms []ProgramInfo

	if err = json.NewDecoder(reader).Decode(&allPrograms); err != nil {
		return nil, err
	}

	err = reader.Close()

	return allPrograms, err
}

// GetProgramDescription returns a set of program descriptions for the given set of program IDs.
func (c Client) GetProgramDescription(programIDs []string) (map[string]ProgramDescription, error) {
	if len(programIDs) > 500 {
		return nil, errors.New("you may only request at most 500 program IDs per request, please lower your request amount")
	}
	url := fmt.Sprint(DefaultBaseURL, APIVersion, "/metadata/description")

	js, jsErr := json.Marshal(programIDs)
	if jsErr != nil {
		return nil, jsErr
	}

	// setup the request
	req, httpErr := http.NewRequest("POST", url, bytes.NewBuffer(js))
	if httpErr != nil {
		return nil, httpErr
	}

	_, data, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	descriptions := make(map[string]ProgramDescription)

	if err = json.Unmarshal(data, &descriptions); err != nil {
		return nil, err
	}

	return descriptions, err
}

// GetLanguageCrossReference returns a map of translated titles and descriptions for the given programIDs.
func (c Client) GetLanguageCrossReference(programIDs []string) (map[string][]LanguageCrossReference, error) {
	// A 500 item limit is not defined in the docs but seems like the reasonable default.
	if len(programIDs) > 500 {
		return nil, errors.New("you may only request at most 500 program IDs per request, please lower your request amount")
	}
	url := fmt.Sprint(DefaultBaseURL, APIVersion, "/xref")

	js, jsErr := json.Marshal(programIDs)
	if jsErr != nil {
		return nil, jsErr
	}

	// setup the request
	req, httpErr := http.NewRequest("POST", url, bytes.NewBuffer(js))
	if httpErr != nil {
		return nil, httpErr
	}

	_, data, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	translations := make(map[string][]LanguageCrossReference)

	if err = json.Unmarshal(data, &translations); err != nil {
		return nil, err
	}

	return translations, err
}

// GetLastModified returns the last modified information for the given station IDs and optional dates.
func (c Client) GetLastModified(requests []StationScheduleRequest) (map[string]map[string]LastModifiedEntry, error) {
	url := fmt.Sprint(DefaultBaseURL, APIVersion, "/schedules/md5")

	s := make(map[string]map[string]LastModifiedEntry)

	js, jsErr := json.Marshal(requests)
	if jsErr != nil {
		return s, jsErr
	}

	req, httpErr := http.NewRequest("POST", url, bytes.NewBuffer(js))
	if httpErr != nil {
		return nil, httpErr
	}

	_, data, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &s)

	return s, err
}

// GetLineups returns a LineupResponse which contains all the lineups subscribed
// to by this account.
func (c Client) GetLineups() (*LineupResponse, error) {
	url := fmt.Sprint(DefaultBaseURL, APIVersion, "/lineups")
	s := new(LineupResponse)

	req, httpErr := http.NewRequest("GET", url, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	_, data, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &s)

	return s, err
}

// DeleteSystemMessage deletes a system message from the status response.
func (c Client) DeleteSystemMessage(messageID string) error {
	url := fmt.Sprint(DefaultBaseURL, "/messages/", messageID)

	req, httpErr := http.NewRequest("DELETE", url, nil)
	if httpErr != nil {
		return httpErr
	}

	_, _, err := c.sendRequest(req)

	return err
}

// GetProgramStillRunning returns the real time status of the given program ID.
func (c Client) GetProgramStillRunning(programID string) (*StillRunningResponse, error) {
	url := fmt.Sprint(DefaultBaseURL, APIVersion, "/metadata/stillRunning/", programID)

	// setup the request
	req, httpErr := http.NewRequest("GET", url, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	_, data, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	stillRunningResp := &StillRunningResponse{}

	if err = json.Unmarshal(data, &stillRunningResp); err != nil {
		return nil, err
	}

	return stillRunningResp, err
}

// GetArtworkForProgramIDs returns artwork for the given programIDs.
func (c Client) GetArtworkForProgramIDs(programIDs []string) ([]ProgramArtworkResponse, error) {
	if len(programIDs) > 500 {
		return nil, errors.New("you may only request at most 500 program IDs per request, please lower your request amount")
	}
	url := fmt.Sprint(DefaultBaseURL, APIVersion, "/metadata/programs")

	js, jsErr := json.Marshal(programIDs)
	if jsErr != nil {
		return nil, jsErr
	}

	// setup the request
	req, httpErr := http.NewRequest("POST", url, bytes.NewBuffer(js))
	if httpErr != nil {
		return nil, httpErr
	}

	_, data, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	programArtwork := make([]ProgramArtworkResponse, 0)

	if err = json.Unmarshal(data, &programArtwork); err != nil {
		return nil, err
	}

	return programArtwork, err
}

// GetArtworkForRootID returns artwork for the given programIDs.
func (c Client) GetArtworkForRootID(rootID string) ([]ProgramArtwork, error) {
	url := fmt.Sprint(DefaultBaseURL, APIVersion, "/metadata/programs/", rootID)

	// setup the request
	req, httpErr := http.NewRequest("GET", url, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	_, data, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	programArtwork := make([]ProgramArtwork, 0)

	if err = json.Unmarshal(data, &programArtwork); err != nil {
		return nil, err
	}

	return programArtwork, err
}

// GetImage returns an image for the given URI.
func (c Client) GetImage(imageURI string) ([]byte, error) {
	url := fmt.Sprint(DefaultBaseURL, APIVersion, "/image/", imageURI)

	if strings.HasPrefix(imageURI, "https://s3.amazonaws.com") {
		url = imageURI
	}

	// setup the request
	req, httpErr := http.NewRequest("GET", url, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	_, data, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetCelebrityArtwork returns artwork for the given programIDs.
func (c Client) GetCelebrityArtwork(celebrityID string) ([]ProgramArtwork, error) {
	url := fmt.Sprint(DefaultBaseURL, APIVersion, "/metadata/celebrity/", celebrityID)

	// setup the request
	req, httpErr := http.NewRequest("GET", url, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	_, data, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	programArtwork := make([]ProgramArtwork, 0)

	if err = json.Unmarshal(data, &programArtwork); err != nil {
		return nil, err
	}

	return programArtwork, err
}

func (c Client) sendRequest(request *http.Request) (response *http.Response, data []byte, err error) {
	request.Header.Set("token", c.Token)

	if request.Method == "POST" {
		request.Header.Set("Content-Type", "application/json")
	}

	response, err = c.HTTP.Do(request)

	if err != nil {
		err = fmt.Errorf("cannot reach server. %v", err)
		return
	}

	if response.StatusCode != http.StatusOK {
		return
	}

	buf := &bytes.Buffer{}
	if _, copyErr := io.Copy(buf, response.Body); copyErr != nil {
		return nil, nil, copyErr
	}

	err = response.Body.Close()
	if err != nil {
		err = fmt.Errorf("cannot read response. %v", err)
	}

	data = buf.Bytes()

	peekBuf := &bytes.Buffer{}
	if _, copyErr := io.Copy(peekBuf, response.Body); copyErr != nil {
		return nil, nil, copyErr
	}

	baseResp := &BaseResponse{}

	if unmarshalErr := json.Unmarshal(peekBuf.Bytes(), baseResp); unmarshalErr == nil {
		if baseResp.Response != "OK" || baseResp.Code != 0 {
			return nil, nil, baseResp
		}
	}

	return
}