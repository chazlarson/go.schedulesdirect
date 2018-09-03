package schedulesdirect

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// A ChangeLineupResponse stores the message after attempting
// to add or delete a lineup.
type ChangeLineupResponse struct {
	*BaseResponse

	ChangesRemaining jsonInt `json:"changesRemaining,omitempty"`
}

// A LineupResponse stores the message after requesting
// to list subscribed lineups.
type LineupResponse struct {
	*BaseResponse

	Lineups []Lineup `json:"lineups,omitempty"`
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
	*BaseResponse

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

// A Station stores a station in a lineup or schedule.
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

// StationPreview is a slim version of Station containing the fields only visible during lineup preview.
type StationPreview struct {
	Affiliate string `json:"affiliate,omitempty"`
	CallSign  string `json:"callsign,omitempty"`
	Channel   string `json:"channel,omitempty"`
	Name      string `json:"name,omitempty"`
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

// AddLineup adds the given lineup uri to the users SchedulesDirect account.
func (c *Client) AddLineup(lineupID string) (*ChangeLineupResponse, error) {
	url := fmt.Sprint(c.BaseURL, APIVersion, "/lineups/", lineupID)

	req, httpErr := http.NewRequest("PUT", url, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	_, data, err := c.SendRequest(req, true)
	if err != nil {
		return nil, err
	}

	r := &ChangeLineupResponse{}

	err = json.Unmarshal(data, &r)
	return r, err
}

// DeleteLineup deletes the given lineup uri from the users SchedulesDirect account.
func (c *Client) DeleteLineup(lineupID string) (*ChangeLineupResponse, error) {
	url := fmt.Sprint(c.BaseURL, APIVersion, "/lineups/", lineupID)

	req, httpErr := http.NewRequest("DELETE", url, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	_, data, err := c.SendRequest(req, true)
	if err != nil {
		return nil, err
	}

	r := &ChangeLineupResponse{}

	err = json.Unmarshal(data, &r)
	return r, err
}

// PreviewLineup returns a slice of StationPreview containing the channels available in the provided lineupID.
func (c *Client) PreviewLineup(lineupID string) ([]StationPreview, error) {
	url := fmt.Sprint(c.BaseURL, APIVersion, "/lineups/preview/", lineupID)

	req, httpErr := http.NewRequest("GET", url, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	_, data, err := c.SendRequest(req, true)
	if err != nil {
		return nil, err
	}

	r := make([]StationPreview, 0)

	err = json.Unmarshal(data, &r)
	return r, err
}

// AutomapLineup accepts the "lineup.json" output as a byte slice from SiliconDust's HDHomerun devices.
// It then runs a comparison against ScheduleDirect's database and returns potential lineup matches.
func (c *Client) AutomapLineup(hdhrLineupJSON []byte) (map[string]int, error) {
	url := fmt.Sprint(c.BaseURL, APIVersion, "/map/lineup")

	req, httpErr := http.NewRequest("POST", url, bytes.NewBuffer(hdhrLineupJSON))
	if httpErr != nil {
		return nil, httpErr
	}

	_, data, err := c.SendRequest(req, true)
	if err != nil {
		return nil, err
	}

	matches := make(map[string]int)

	err = json.Unmarshal(data, &matches)
	return matches, err
}

// SubmitLineup should be called if AutomapLineup doesn't return candidates after you identify
// the lineup you were trying to find via automapping.
func (c *Client) SubmitLineup(hdhrLineupJSON []byte, lineupID string) error {
	url := fmt.Sprint(c.BaseURL, APIVersion, "/map/lineup/", lineupID)

	req, httpErr := http.NewRequest("POST", url, bytes.NewBuffer(hdhrLineupJSON))
	if httpErr != nil {
		return httpErr
	}

	_, _, err := c.SendRequest(req, true)
	return err
}

// GetHeadends returns the map of headends for the given country and postal code.
func (c *Client) GetHeadends(countryCode, postalCode string) ([]Headend, error) {
	params := url.Values{}
	params.Add("country", countryCode)
	params.Add("postalcode", postalCode)
	uriPart := fmt.Sprint("/headends?", params.Encode())
	url := fmt.Sprint(c.BaseURL, APIVersion, uriPart)

	req, httpErr := http.NewRequest("GET", url, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	_, data, err := c.SendRequest(req, true)
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
func (c *Client) GetChannels(lineupID string, verbose bool) (*ChannelResponse, error) {
	url := fmt.Sprint(c.BaseURL, APIVersion, "/lineups/", lineupID)

	req, httpErr := http.NewRequest("GET", url, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	if verbose {
		req.Header.Add("verboseMap", "true")
	}

	_, data, err := c.SendRequest(req, true)
	if err != nil {
		return nil, err
	}

	h := new(ChannelResponse)

	err = json.Unmarshal(data, &h)
	return h, err
}

// GetLineups returns a LineupResponse which contains all the lineups subscribed
// to by this account.
func (c *Client) GetLineups() (*LineupResponse, error) {
	url := fmt.Sprint(c.BaseURL, APIVersion, "/lineups")
	s := new(LineupResponse)

	req, httpErr := http.NewRequest("GET", url, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	_, data, err := c.SendRequest(req, true)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &s)

	return s, err
}
