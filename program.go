package schedulesdirect

import (
	"encoding/json"
	"time"
)

// A SyndicationType stores syndication information for a program
type SyndicationType struct {
	Source string `json:"source"`
	Type   string `json:"type"`
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
	ProgramID           string          `json:"programID,omitempty"`
	AirDateTime         time.Time       `json:"airDateTime,omitempty"`
	MD5                 string          `json:"md5,omitempty"`
	Duration            int             `json:"duration,omitempty"`
	LiveTapeDelay       string          `json:"liveTapeDelay,omitempty"`
	IsPremiereOrFinale  PremiereType    `json:"isPremiereOrFinale"`
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
	Ratings             []ContentRating `json:"ratings,omitempty"`
	ProgramPart         Part            `json:"multipart,omitempty"`
	VideoProperties     []string        `json:"videoProperties,omitempty"`
}

// A ProgramInfo type stores information for a program.
type ProgramInfo struct {
	BaseResponse BaseResponse

	Animation         string                   `json:"animation"`
	Audience          string                   `json:"audience"`
	Awards            []Award                  `json:"awards"`
	Cast              []Person                 `json:"cast"`
	ContentAdvisory   []string                 `json:"contentAdvisory"`
	ContentRating     []ContentRating          `json:"contentRating"`
	Crew              []Person                 `json:"crew"`
	Descriptions      map[string][]Description `json:"descriptions"`
	Duration          int                      `json:"duration"`
	EntityType        string                   `json:"entityType"`
	EpisodeTitle150   string                   `json:"episodeTitle150"`
	EventDetails      EventDetails             `json:"eventDetails"`
	Genres            []string                 `json:"genres"`
	HasEpisodeArtwork bool                     `json:"hasEpisodeArtwork"`
	HasImageArtwork   bool                     `json:"hasImageArtwork"`
	HasMovieArtwork   bool                     `json:"hasMovieArtwork"`
	HasSeriesArtwork  bool                     `json:"hasSeriesArtwork"`
	HasSportsArtwork  bool                     `json:"hasSportsArtwork"`
	Holiday           string                   `json:"holiday"`
	Keywords          map[string][]string      `json:"keyWords"`
	MD5               string                   `json:"md5"`
	Metadata          []map[string]Metadata    `json:"metadata"`
	Movie             Movie                    `json:"movie"`
	OfficialURL       string                   `json:"officialURL"`
	OriginalAirDate   Date                     `json:"originalAirDate"`
	ProgramID         string                   `json:"programID"`
	Recommendations   []Recommendation         `json:"recommendations"`
	ResourceID        string                   `json:"resourceID"`
	ShowType          string                   `json:"showType"`
	Titles            []Title                  `json:"titles"`
}

// HasArtwork returns true if the Program has artwork available.
func (p *ProgramInfo) HasArtwork() bool {
	return p.HasEpisodeArtwork || p.HasImageArtwork || p.HasMovieArtwork || p.HasSeriesArtwork || p.HasSportsArtwork
}

// Award is a award given to a program.
type Award struct {
	AwardName string `json:"awardName"`
	Category  string `json:"category"`
	Name      string `json:"name"`
	PersonID  string `json:"personId"`
	Recipient string `json:"recipient"`
	Won       bool   `json:"won"`
	Year      Date   `json:"year"`
}

// Person stores information for an acting credit or crew member.
type Person struct {
	PersonID      string `json:"personId,omitmepty"`
	NameID        string `json:"nameId,omitempty"`
	Name          string `json:"name,omitempty"`
	Role          string `json:"role,omitempty"`
	CharacterName string `json:"characterName,omitempty"`
	BillingOrder  string `json:"billingOrder,omitempty"`
}

// A ContentRating stores ratings board information for a program
type ContentRating struct {
	Body    string `json:"body"`
	Code    string `json:"code"`
	Country string `json:"country"`
}

// Description provides a generic description of a program.
type Description struct {
	Description string `json:"description"`
	Language    string `json:"descriptionLanguage"`
}

// A Movie type stores information about a movie
type Movie struct {
	Duration      int                  `json:"duration"`
	QualityRating []MovieQualityRating `json:"qualityRating"`
	Year          Date                 `json:"year"`
}

// Metadata stores meta information for a program.
type Metadata struct {
	Episode       int `json:"episode"`
	EpisodeID     int `json:"episodeID"`
	Season        int `json:"season"`
	SeriesID      int `json:"seriesID"`
	TotalEpisodes int `json:"totalEpisodes"`
	TotalSeasons  int `json:"totalSeasons"`
}

// EventDetails contains details about the sporting program related to a game.
type EventDetails struct {
	GameDate Date   `json:"gameDate"`
	Teams    []Team `json:"teams"`
	Venue    string `json:"venue100"`
}

// A MovieQualityRating describes ratings for the quality of a movie.
type MovieQualityRating struct {
	Increment   string `json:"increment"`
	MaxRating   string `json:"maxRating"`
	MinRating   string `json:"minRating"`
	Rating      string `json:"rating"`
	RatingsBody string `json:"ratingsBody"`
}

// Team is a sports team that participated in a game program.
type Team struct {
	IsHome bool   `json:"isHome"`
	Name   string `json:"name"`
	Score  string `json:"score"`
}

// Recommendation is a related content recommendation.
type Recommendation struct {
	ProgramID string `json:"programID"`
	Title120  string `json:"title120"`
}

// Title contains the title of a program.
type Title struct {
	Title120 string `json:"title120"`
}

// Part stores the information for a part
type Part struct {
	PartNumber int `json:"partNumber"`
	TotalParts int `json:"totalParts"`
}

// ProgramDescription provides a generic description of a program.
type ProgramDescription struct {
	Code            int    `json:"code"`
	Description100  string `json:"description100"`
	Description1000 string `json:"description1000"`
}

// ProgramArtwork describes a single piece of artwork related to a program.
type ProgramArtwork struct {
	Aspect   string            `json:"aspect"`
	Category string            `json:"category"`
	Height   int               `json:"height,string"`
	Primary  string            `json:"primary"`
	Size     string            `json:"size"`
	Text     string            `json:"text"`
	Tier     string            `json:"tier"`
	URI      string            `json:"uri"`
	Width    int               `json:"width,string"`
	Caption  map[string]string `json:"caption"`
}

// ProgramArtworkResponse is a container struct for artwork relating to a program.
type ProgramArtworkResponse struct {
	ProgramID string            `json:"programID"`
	Error     *BaseResponse     `json:"-"`
	Artwork   *[]ProgramArtwork `json:"-"`
	wrapper   struct {
		PID  string          `json:"programID"`
		Data json.RawMessage `json:"data"`
	}
}

// UnmarshalJSON unmarshals the JSON into the ProgramArtworkResponse.
func (ar *ProgramArtworkResponse) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &ar.wrapper); err != nil {
		return err
	}
	ar.ProgramID = ar.wrapper.PID
	if ar.wrapper.Data[0] == '[' {
		return json.Unmarshal(ar.wrapper.Data, &ar.Artwork)
	}
	return json.Unmarshal(ar.wrapper.Data, &ar.Error)
}
