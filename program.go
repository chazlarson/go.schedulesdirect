package schedulesdirect

import (
	"encoding/json"
	"fmt"
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

// ArtworkTier is the "level" of the artwork
type ArtworkTier string

const (
	// SeriesTier - Series - image represents of series, regardless of season (banner, iconic, staple, cast, logo).
	SeriesTier ArtworkTier = "Series"
	// SeasonTier - Season - image represents specific season of series (banner, iconic, cast, logo).
	SeasonTier ArtworkTier = "Season"
	// EpisodeTier - Episode - image represents specific episode of series (iconics only).
	EpisodeTier ArtworkTier = "Episode"
	// TeamEventNewTier - Team Event (New) - image represents team vs. team events (banner).
	TeamEventNewTier ArtworkTier = "Team Event (New)"
	// OrganizationTier - Organization - image represents organization associated with sport (logo).
	OrganizationTier ArtworkTier = "Organization"
	// ConferenceTier - Conference - image represents conference associated with team (logo).
	ConferenceTier ArtworkTier = "Conference"
	// SportTier - Sport - image represents the sport associated with a program (banner, iconic).
	SportTier ArtworkTier = "Sport"
	// SportEventTier - Sport Event -  image represents the sport event associated with a program (banner, iconic).
	SportEventTier ArtworkTier = "Sport Event"
	// CollegeTier - College - image represents the college associated with a program (logo).
	CollegeTier ArtworkTier = "College"
	// TeamTier - Team - image represents the team associated with a program (logo).
	TeamTier ArtworkTier = "Team"
)

// ArtworkCategory describes the kind of artwork.
type ArtworkCategory string

const (
	// Banner is a source-provided image, usually shows cast ensemble with source-provided text
	Banner ArtworkCategory = "Banner"
	// BannerL1 is the same as Banner
	BannerL1 ArtworkCategory = "Banner-L1"
	// BannerL1T is the same as banner with text
	BannerL1T ArtworkCategory = "Banner-L1T"
	// BannerL2 is a source-provided image with plain text
	BannerL2 ArtworkCategory = "Banner-L2"
	// BannerL3 is a stock photo image with plain text
	BannerL3 ArtworkCategory = "Banner-L3"
	// BannerLO is a banner with Logo Only
	BannerLO ArtworkCategory = "Banner-LO"
	// BannerLOT is a banner with Logo Only + Text indicating season number
	BannerLOT ArtworkCategory = "Banner-LOT"
	// Iconic is a representative series/season/episode image, no text
	Iconic ArtworkCategory = "Iconic"
	// Staple is a the staple image is intended to cover programs which do not have a unique banner image.
	Staple ArtworkCategory = "Staple"
	// CastEnsemble is a cast ensemble, no text
	CastEnsemble ArtworkCategory = "Cast Ensemble"
	// CastInCharacter is a individual cast member, no text
	CastInCharacter ArtworkCategory = "Cast in Character"
	// Logo is a official logo for program, sports organization, sports conference, or TV station
	Logo ArtworkCategory = "Logo"
	// BoxArt is a DVD box art, for movies only
	BoxArt ArtworkCategory = "Box Art"
	// PosterArt is a theatrical movie poster, standard sizes
	PosterArt ArtworkCategory = "Poster Art"
	// SceneStill is a movie photo, legacy sizes
	SceneStill ArtworkCategory = "Scene Still"
	// Photo is the same as Scene Still
	Photo ArtworkCategory = "Photo"
	// PhotoHeadshot is a celebrity image
	PhotoHeadshot ArtworkCategory = "Photo-headshot"
	// VODArt is a image used for video on demand.
	VODArt ArtworkCategory = "VOD Art"
)

// ArtworkSize is the size class of the Artwork.
type ArtworkSize string

const (
	// ExtraSmallArtworkSize means the artwork size is extra small (Xs).
	ExtraSmallArtworkSize ArtworkSize = "Xs"
	// SmallArtworkSize means the artwork size is small (Sm).
	SmallArtworkSize ArtworkSize = "Sm"
	// MediumArtworkSize means the artwork size is medium (Md).
	MediumArtworkSize ArtworkSize = "Md"
	// LargeArtworkSize means the artwork size is large (Lg).
	LargeArtworkSize ArtworkSize = "Lg"
	// MasterArtworkSize means the artwork size is master (Ms).
	MasterArtworkSize ArtworkSize = "Ms"
)

// ArtworkAspectRatio is the aspect ratio of the Artwork.
type ArtworkAspectRatio string

const (
	// SixteenByNineAspectRatio means the artwork has a aspect ratio of 16 by 9.
	SixteenByNineAspectRatio ArtworkAspectRatio = "16x9"
	// FourByThreeAspectRatio means the artwork has a aspect ratio of 4 by 3.
	FourByThreeAspectRatio ArtworkAspectRatio = "4x3"
	// ThreeByFourAspectRatio means the artwork has a aspect ratio of 3 by 4.
	ThreeByFourAspectRatio ArtworkAspectRatio = "3x4"
	// TwoByThreeAspectRatio means the artwork has a aspect ratio of 2 by 3.
	TwoByThreeAspectRatio ArtworkAspectRatio = "2x3"
	// OneByOneAspectRatio means the artwork has a aspect ratio of 1 by 1.
	OneByOneAspectRatio ArtworkAspectRatio = "1x1"
)

// ArtworkCaption is the caption assigned to an Artwork.
type ArtworkCaption struct {
	Content  string `json:"content"`
	Language string `json:"lang"`
}

// Artwork describes a single piece of artwork related to a program.
type Artwork struct {
	Aspect   ArtworkAspectRatio `json:"aspect,omitempty"`
	Category ArtworkCategory    `json:"category,omitempty"`
	Height   int                `json:"height,string,omitempty"`
	Primary  ConvertibleBoolean `json:"primary,omitempty"`
	Size     ArtworkSize        `json:"size,omitempty"`
	Text     ConvertibleBoolean `json:"text,omitempty"`
	Tier     ArtworkTier        `json:"tier,omitempty"`
	URI      string             `json:"uri,omitempty"`
	Width    int                `json:"width,string,omitempty"`
	Caption  ArtworkCaption     `json:"caption,omitempty"`
}

// ArtworkResponse is a container struct for artwork relating to a program.
type ArtworkResponse struct {
	ProgramID string        `json:"programID,omitempty"`
	Error     *BaseResponse `json:"-"`
	Artwork   *[]Artwork    `json:"-"`
	wrapper   struct {
		PID  string          `json:"programID,omitempty"`
		Data json.RawMessage `json:"data,omitempty"`
	}
}

// UnmarshalJSON unmarshals the JSON into the ArtworkResponse.
func (ar *ArtworkResponse) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &ar.wrapper); err != nil {
		return err
	}
	ar.ProgramID = ar.wrapper.PID
	if ar.wrapper.Data[0] == '[' {
		return json.Unmarshal(ar.wrapper.Data, &ar.Artwork)
	}
	return json.Unmarshal(ar.wrapper.Data, &ar.Error)
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
