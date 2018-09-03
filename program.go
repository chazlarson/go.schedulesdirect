package schedulesdirect

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// A ProgramInfo type stores information for a program.
type ProgramInfo struct {
	BaseResponse *BaseResponse `json:"-,omitempty"`

	Animation         Animation                `json:"animation,omitempty"`
	Audience          Audience                 `json:"audience,omitempty"`
	Awards            []Award                  `json:"awards,omitempty"`
	Cast              []Person                 `json:"cast,omitempty"`
	ContentAdvisory   []string                 `json:"contentAdvisory,omitempty"`
	ContentRating     []ContentRating          `json:"contentRating,omitempty"`
	Crew              []Person                 `json:"crew,omitempty"`
	Descriptions      map[string][]Description `json:"descriptions,omitempty"`
	Duration          int                      `json:"duration,omitempty"`
	EntityType        EntityType               `json:"entityType,omitempty"`
	EpisodeImage      *Artwork                 `json:"episodeImage,omitempty"`
	EpisodeTitle150   string                   `json:"episodeTitle150,omitempty"`
	EventDetails      *EventDetails            `json:"eventDetails,omitempty"`
	Genres            []string                 `json:"genres,omitempty"`
	HasEpisodeArtwork bool                     `json:"hasEpisodeArtwork,omitempty"`
	HasImageArtwork   bool                     `json:"hasImageArtwork,omitempty"`
	HasMovieArtwork   bool                     `json:"hasMovieArtwork,omitempty"`
	HasSeriesArtwork  bool                     `json:"hasSeriesArtwork,omitempty"`
	HasSportsArtwork  bool                     `json:"hasSportsArtwork,omitempty"`
	Holiday           string                   `json:"holiday,omitempty"`
	Keywords          map[string][]string      `json:"keyWords,omitempty"`
	MD5               string                   `json:"md5,omitempty"`
	Metadata          []map[string]Metadata    `json:"metadata,omitempty"`
	Movie             *Movie                   `json:"movie,omitempty"`
	OfficialURL       string                   `json:"officialURL,omitempty"`
	OriginalAirDate   *Date                    `json:"originalAirDate,omitempty"`
	ProgramID         string                   `json:"programID,omitempty"`
	Recommendations   []Recommendation         `json:"recommendations,omitempty"`
	ResourceID        string                   `json:"resourceID,omitempty"`
	ShowType          ShowSubType              `json:"showType,omitempty"`
	Titles            []Title                  `json:"titles,omitempty"`
}

// HasArtwork returns true if the Program has artwork available.
func (p *ProgramInfo) HasArtwork() bool {
	return p.HasEpisodeArtwork || p.HasImageArtwork || p.HasMovieArtwork || p.HasSeriesArtwork || p.HasSportsArtwork
}

// ShowID is just a helper wrapper around GetShowIDForEpisodeID.
func (p *ProgramInfo) ShowID() string {
	return GetShowIDForEpisodeID(p.ProgramID)
}

// ArtworkLookupIDs returns a string slice of IDs that can be used to look up artwork for the program by.
//
// This effectively returns a slice containing the programID while also appending the
// SH program ID if the programID begins with EP.
func (p *ProgramInfo) ArtworkLookupIDs() []string {
	if p.HasEpisodeArtwork && p.ShowID() != "" { // If the program has episode artwork and a SH ID (e.g. EP024874280035)
		return []string{p.ProgramID, p.ShowID()} // return []string{"EP024874280035", "SH024874280000"}
	} else if !p.HasEpisodeArtwork && p.ProgramID[0:2] == "EP" { // If the program doesn't have episode artwork but is an episode (e.g. EP027100890371)
		return []string{fmt.Sprintf("SH%s0000", p.ProgramID[2:10])} // return []string{"SH027100890000"}
	}
	return []string{p.ProgramID}
}

// Animation is the type of animation employed by the Program
type Animation string

const (
	// Animated means the program is animated.
	Animated Animation = "Animated"
	// Anime means the program uses the Anime style of animation.
	Anime Animation = "Anime"
	// LiveActionAnimated means the program uses a combination of live action and animated sequences.
	LiveActionAnimated Animation = "Live action/animated"
	// LiveActionAnime means the program uses a combination of live action and anime sequences.
	LiveActionAnime Animation = "Live action/anime"
)

// Audience indicates program target audience, derived from genres.
type Audience string

const (
	// ChildrensAudience means the program is suitable for children only.
	ChildrensAudience Audience = "Children"
	// AdultsOnlyAudience means the program is suitable for adults only.
	AdultsOnlyAudience Audience = "Adults only"
)

// EntityType is the program type
type EntityType string

const (
	// EpisodeEntityType means the program is a episode of a TV show.
	EpisodeEntityType EntityType = "Episode"
	// MovieEntityType means the program is a movie.
	MovieEntityType EntityType = "Movie"
	// ShowEntityType means the program is a TV show.
	ShowEntityType EntityType = "Show"
	// SportsEntityType means the program is a sporting event or related program.
	SportsEntityType EntityType = "Sports"
)

// ShowSubType is the program subtype.
type ShowSubType string

const (
	// FeatureFilm means the program sub type is a feature film.
	FeatureFilm ShowSubType = "Feature Film"
	// MiniSeries means the program sub type is a miniseries.
	MiniSeries ShowSubType = "Miniseries"
	// PaidProgramming means the program sub type is a paid programming.
	PaidProgramming ShowSubType = "Paid Programming"
	// Series means the program sub type is a series.
	Series ShowSubType = "Series"
	// ShortFilm means the program sub type is a short film.
	ShortFilm ShowSubType = "Short Film"
	// Special means the program sub type is a special.
	Special ShowSubType = "Special"
	// SportsEvent means the program sub type is a sports event.
	SportsEvent ShowSubType = "Sports event"
	// SportsNonEvent means the program sub type is a sports non-event.
	SportsNonEvent ShowSubType = "Sports non-event"
	// TheatreEvent means the program sub type is a theatre event.
	TheatreEvent ShowSubType = "Theatre Event"
	// TVMovie means the program sub type is a TV movie.
	TVMovie ShowSubType = "TV Movie"
)

// Award is a award given to a program.
type Award struct {
	AwardName string `json:"awardName,omitempty"`
	Category  string `json:"category,omitempty"`
	Name      string `json:"name,omitempty"`
	PersonID  string `json:"personId,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Won       bool   `json:"won,omitempty"`
	Year      *Date  `json:"year,omitempty"`
}

// Person stores information for an acting credit or crew member.
type Person struct {
	PersonID      string `json:"personId,omitmepty,omitempty"`
	NameID        string `json:"nameId,omitempty"`
	Name          string `json:"name,omitempty"`
	Role          string `json:"role,omitempty"`
	CharacterName string `json:"characterName,omitempty"`
	BillingOrder  string `json:"billingOrder,omitempty"`
}

// A ContentRating stores ratings board information for a program
type ContentRating struct {
	Body    string `json:"body,omitempty"`
	Code    string `json:"code,omitempty"`
	Country string `json:"country,omitempty"`
}

// Description provides a generic description of a program.
type Description struct {
	Description string `json:"description,omitempty"`
	Language    string `json:"descriptionLanguage,omitempty"`
}

// A Movie type stores information about a movie
type Movie struct {
	Duration      int                  `json:"duration,omitempty"`
	QualityRating []MovieQualityRating `json:"qualityRating,omitempty"`
	Year          *Date                `json:"year,omitempty"`
}

// Metadata stores meta information for a program.
type Metadata struct {
	Episode       int `json:"episode,omitempty"`
	EpisodeID     int `json:"episodeID,omitempty"`
	Season        int `json:"season,omitempty"`
	SeriesID      int `json:"seriesID,omitempty"`
	TotalEpisodes int `json:"totalEpisodes,omitempty"`
	TotalSeasons  int `json:"totalSeasons,omitempty"`
}

// EventDetails contains details about the sporting program related to a game.
type EventDetails struct {
	GameDate *Date  `json:"gameDate,omitempty"`
	Teams    []Team `json:"teams,omitempty"`
	Venue    string `json:"venue100,omitempty"`
	SubType  string `json:"subType,omitempty"`
}

// A MovieQualityRating describes ratings for the quality of a movie.
type MovieQualityRating struct {
	Increment   string `json:"increment,omitempty"`
	MaxRating   string `json:"maxRating,omitempty"`
	MinRating   string `json:"minRating,omitempty"`
	Rating      string `json:"rating,omitempty"`
	RatingsBody string `json:"ratingsBody,omitempty"`
}

// Team is a sports team that participated in a game program.
type Team struct {
	IsHome bool   `json:"isHome,omitempty"`
	Name   string `json:"name,omitempty"`
	Score  string `json:"score,omitempty"`
}

// Recommendation is a related content recommendation.
type Recommendation struct {
	ProgramID string `json:"programID,omitempty"`
	Title120  string `json:"title120,omitempty"`
}

// Title contains the title of a program.
type Title struct {
	Title120 string `json:"title120,omitempty"`
}

// Part stores the information for a part
type Part struct {
	PartNumber int `json:"partNumber,omitempty"`
	TotalParts int `json:"totalParts,omitempty"`
}

// ProgramDescription provides a generic description of a program.
type ProgramDescription struct {
	Code            int    `json:"code,omitempty"`
	Description100  string `json:"description100,omitempty"`
	Description1000 string `json:"description1000,omitempty"`
}

// LanguageCrossReference provides translated titles and descriptions for a program.
type LanguageCrossReference struct {
	*BaseResponse

	DescriptionLanguage     string `json:"descriptionLanguage,omitempty"`
	DescriptionLanguageName string `json:"descriptionLanguageName,omitempty"`
	MD5                     string `json:"md5,omitempty"`
	ProgramID               string `json:"programID,omitempty"`
	TitleLanguage           string `json:"titleLanguage,omitempty"`
	TitleLanguageName       string `json:"titleLanguageName,omitempty"`
}

// A StillRunningResponse describes the current real time state of a program.
type StillRunningResponse struct {
	*BaseResponse

	EventStartDateTime *time.Time `json:"eventStartDateTime,omitempty"`
	IsComplete         bool       `json:"isComplete,omitempty"`
	ProgramID          string     `json:"programID,omitempty"`
	Result             struct {
		AwayTeam *Team `json:"awayTeam,omitempty"`
		HomeTeam *Team `json:"homeTeam,omitempty"`
	} `json:"result,omitempty"`
}

// GetShowIDForEpisodeID returns a string containing the a SH program ID if
// the input program has an ID beginning with EP.
//
// If programID is not a EP program ID it will return an empty string.
func GetShowIDForEpisodeID(programID string) string {
	if programID[0:2] == "EP" {
		return fmt.Sprintf("SH%s%s", programID[2:10], "0000")
	}
	return ""
}

// GetProgramInfo returns the set of program details for the given set of programs.
//
// If more than 5000 Program IDs are provided, the client will automatically
// chunk the slice into groups of 5000 IDs and return all responses to you.
func (c *Client) GetProgramInfo(programIDs []string) ([]ProgramInfo, error) {
	// If user passed more than 5000 programIDs, let's help them out by
	// chunking the requests for them.
	// Obviously you can disable this behavior by passing less than 5000 IDs.
	if len(programIDs) > 5000 {
		allResponses := make([]ProgramInfo, len(programIDs))
		for _, chunk := range chunkStringSlice(programIDs, 5000) {
			resp, err := c.GetProgramInfo(chunk)
			if err != nil {
				return nil, err
			}
			allResponses = append(allResponses, resp...)
		}
		return allResponses, nil
	}

	url := fmt.Sprint(c.BaseURL, APIVersion, "/programs")

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

	_, data, err := c.SendRequest(req, true)
	if err != nil {
		return nil, err
	}

	// create the programs slice
	allPrograms := make([]ProgramInfo, 0)

	if err = json.Unmarshal(data, &allPrograms); err != nil {
		return nil, err
	}

	return allPrograms, err
}

// GetProgramDescription returns a set of program descriptions for the given set of program IDs.
//
// If more than 500 Program IDs are provided, the client will automatically
// chunk the slice into groups of 500 IDs and return all responses to you.
func (c *Client) GetProgramDescription(programIDs []string) (map[string]ProgramDescription, error) {
	// If user passed more than 500 programIDs, let's help them out by
	// chunking the requests for them.
	// Obviously you can disable this behavior by passing less than 500 IDs.
	if len(programIDs) > 500 {
		allResponses := make(map[string]ProgramDescription)
		for _, chunk := range chunkStringSlice(programIDs, 500) {
			resp, err := c.GetProgramDescription(chunk)
			if err != nil {
				return nil, err
			}
			for key, val := range resp {
				allResponses[key] = val
			}
		}
		return allResponses, nil
	}

	url := fmt.Sprint(c.BaseURL, APIVersion, "/metadata/description")

	js, jsErr := json.Marshal(programIDs)
	if jsErr != nil {
		return nil, jsErr
	}

	// setup the request
	req, httpErr := http.NewRequest("POST", url, bytes.NewBuffer(js))
	if httpErr != nil {
		return nil, httpErr
	}

	_, data, err := c.SendRequest(req, true)
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
//
// If more than 500 Program IDs are provided, the client will automatically
// chunk the slice into groups of 500 IDs and return all responses to you.
func (c *Client) GetLanguageCrossReference(programIDs []string) (map[string][]LanguageCrossReference, error) {
	// A 500 item limit is not defined in the docs but seems like the reasonable default.
	// If user passed more than 500 programIDs, let's help them out by
	// chunking the requests for them.
	// Obviously you can disable this behavior by passing less than 500 IDs.
	if len(programIDs) > 500 {
		allResponses := make(map[string][]LanguageCrossReference)
		for _, chunk := range chunkStringSlice(programIDs, 500) {
			resp, err := c.GetLanguageCrossReference(chunk)
			if err != nil {
				return nil, err
			}
			for key, val := range resp {
				allResponses[key] = val
			}
		}
		return allResponses, nil
	}

	url := fmt.Sprint(c.BaseURL, APIVersion, "/xref")

	js, jsErr := json.Marshal(programIDs)
	if jsErr != nil {
		return nil, jsErr
	}

	// setup the request
	req, httpErr := http.NewRequest("POST", url, bytes.NewBuffer(js))
	if httpErr != nil {
		return nil, httpErr
	}

	_, data, err := c.SendRequest(req, true)
	if err != nil {
		return nil, err
	}

	translations := make(map[string][]LanguageCrossReference)

	if err = json.Unmarshal(data, &translations); err != nil {
		return nil, err
	}

	return translations, err
}

// GetProgramStillRunning returns the real time status of the given program ID.
func (c *Client) GetProgramStillRunning(programID string) (*StillRunningResponse, error) {
	url := fmt.Sprint(c.BaseURL, APIVersion, "/metadata/stillRunning/", programID)

	// setup the request
	req, httpErr := http.NewRequest("GET", url, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	_, data, err := c.SendRequest(req, true)
	if err != nil {
		return nil, err
	}

	stillRunningResp := &StillRunningResponse{}

	if err = json.Unmarshal(data, &stillRunningResp); err != nil {
		return nil, err
	}

	return stillRunningResp, err
}
