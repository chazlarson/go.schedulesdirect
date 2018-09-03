package schedulesdirect

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// A SyndicationType stores syndication information for a program
type SyndicationType struct {
	Source string `json:"source,omitempty"`
	Type   string `json:"type,omitempty"`
}

// PremiereType is used for enumerating the IsPremiereOrFinale field of a Program.
type PremiereType string

const (
	// Finale denotes the end of a show.
	Finale PremiereType = "Finale"
	// Premiere is the beginning of a show.
	Premiere PremiereType = "Premiere"
	// SeasonFinale is the end of a season of a show.
	SeasonFinale PremiereType = "Season Finale"
	// SeasonPremiere is the beginning of a season of a show.
	SeasonPremiere PremiereType = "Season Premiere"
	// SeriesFinale is the end of a show.
	SeriesFinale PremiereType = "Series Finale"
	// SeriesPremiere is the beginning of a show.
	SeriesPremiere PremiereType = "Series Premiere"
)

// LiveTapeDelay signifies if this showing is Live, or Tape Delayed.
type LiveTapeDelay string

const (
	// Live means the program is being shown in real time.
	Live LiveTapeDelay = "Live"
	// Tape means the program was previously recorded.
	Tape LiveTapeDelay = "Tape"
	// Delayed means the program is being intentionally delayed ("broadcast delay").
	Delayed LiveTapeDelay = "Delayed"
)

// A Program stores the information to describing a single television program.
type Program struct {
	ProgramID           string           `json:"programID,omitempty"`
	AirDateTime         *time.Time       `json:"airDateTime,omitempty"`
	MD5                 string           `json:"md5,omitempty"`
	Duration            int              `json:"duration,omitempty"`
	LiveTapeDelay       string           `json:"liveTapeDelay,omitempty"`
	IsPremiereOrFinale  *PremiereType    `json:"isPremiereOrFinale,omitempty"`
	New                 bool             `json:"new,omitempty"`
	CableInTheClassroom bool             `json:"cableInTheClassRoom,omitempty"`
	Catchup             bool             `json:"catchup,omitempty"`   // - typically only found outside of North America
	Continued           bool             `json:"continued,omitempty"` // - typically only found outside of North America
	Education           bool             `json:"educational,omitempty"`
	JoinedInProgress    bool             `json:"joinedInProgress,omitempty"`
	LeftInProgress      bool             `json:"leftInProgress,omitempty"`
	Premiere            bool             `json:"premiere,omitempty"`          //- Should only be found in Miniseries and Movie program types.
	ProgramBreak        bool             `json:"programBreak,omitempty"`      // - Program stops and will restart later (frequently followed by a continued). Typically only found outside of North America.
	Repeat              bool             `json:"repeat,omitempty"`            // - An encore presentation. Repeat should only be found on a second telecast of sporting events.
	Signed              bool             `json:"signed,omitempty"`            //- Program has an on-screen person providing sign-language translation.
	SubjectToBlackout   bool             `json:"subjectToBlackout,omitempty"` //subjectToBlackout
	TimeApproximate     bool             `json:"timeApproximate,omitempty"`
	AudioProperties     []string         `json:"audioProperties,omitempty"`
	Syndication         *SyndicationType `json:"syndication,omitempty"`
	Ratings             []ContentRating  `json:"ratings,omitempty"`
	ProgramPart         *Part            `json:"multipart,omitempty"`
	VideoProperties     []string         `json:"videoProperties,omitempty"`
}

// ShowID is just a helper wrapper around GetShowIDForEpisodeID.
func (p *Program) ShowID() string {
	return GetShowIDForEpisodeID(p.ProgramID)
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

// GetSchedules returns the set of schedules requested.  As a whole the response is not valid json but each individual line is valid.
func (c *Client) GetSchedules(requests []StationScheduleRequest) ([]Schedule, error) {
	url := fmt.Sprint(c.BaseURL, APIVersion, "/schedules")

	js, jsErr := json.Marshal(requests)
	if jsErr != nil {
		return nil, jsErr
	}

	//setup the request
	req, httpErr := http.NewRequest("POST", url, bytes.NewBuffer(js))
	if httpErr != nil {
		return nil, httpErr
	}

	_, data, err := c.SendRequest(req, true)
	if err != nil {
		return nil, err
	}

	var h []Schedule

	err = json.Unmarshal(data, &h)
	return h, err
}

// GetLastModified returns the last modified information for the given station IDs and optional dates.
func (c *Client) GetLastModified(requests []StationScheduleRequest) (map[string]map[string]LastModifiedEntry, error) {
	url := fmt.Sprint(c.BaseURL, APIVersion, "/schedules/md5")

	s := make(map[string]map[string]LastModifiedEntry)

	js, jsErr := json.Marshal(requests)
	if jsErr != nil {
		return s, jsErr
	}

	req, httpErr := http.NewRequest("POST", url, bytes.NewBuffer(js))
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
