package schedulesdirect

import (
	"encoding/json"
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

// A ProgramInfo type stores information for a program.
type ProgramInfo struct {
	BaseResponse *BaseResponse `json:"-,omitempty"`

	Animation         string                   `json:"animation,omitempty"`
	Audience          string                   `json:"audience,omitempty"`
	Awards            []Award                  `json:"awards,omitempty"`
	Cast              []Person                 `json:"cast,omitempty"`
	ContentAdvisory   []string                 `json:"contentAdvisory,omitempty"`
	ContentRating     []ContentRating          `json:"contentRating,omitempty"`
	Crew              []Person                 `json:"crew,omitempty"`
	Descriptions      map[string][]Description `json:"descriptions,omitempty"`
	Duration          int                      `json:"duration,omitempty"`
	EntityType        string                   `json:"entityType,omitempty"`
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
	ShowType          string                   `json:"showType,omitempty"`
	Titles            []Title                  `json:"titles,omitempty"`
}

// HasArtwork returns true if the Program has artwork available.
func (p *ProgramInfo) HasArtwork() bool {
	return p.HasEpisodeArtwork || p.HasImageArtwork || p.HasMovieArtwork || p.HasSeriesArtwork || p.HasSportsArtwork
}

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

// ProgramArtwork describes a single piece of artwork related to a program.
type ProgramArtwork struct {
	Aspect   string            `json:"aspect,omitempty"`
	Category string            `json:"category,omitempty"`
	Height   int               `json:"height,string,omitempty"`
	Primary  string            `json:"primary,omitempty"`
	Size     string            `json:"size,omitempty"`
	Text     string            `json:"text,omitempty"`
	Tier     string            `json:"tier,omitempty"`
	URI      string            `json:"uri,omitempty"`
	Width    int               `json:"width,string,omitempty"`
	Caption  map[string]string `json:"caption,omitempty"`
}

// ProgramArtworkResponse is a container struct for artwork relating to a program.
type ProgramArtworkResponse struct {
	ProgramID string            `json:"programID,omitempty"`
	Error     *BaseResponse     `json:"-,omitempty"`
	Artwork   *[]ProgramArtwork `json:"-,omitempty"`
	wrapper   struct {
		PID  string          `json:"programID,omitempty"`
		Data json.RawMessage `json:"data,omitempty"`
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
