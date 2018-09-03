package schedulesdirect

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

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

// GetArtworkForProgramIDs returns artwork for the given programIDs.
//
// If more than 500 Program IDs are provided, the client will automatically
// chunk the slice into groups of 500 IDs and return all responses to you.
func (c *Client) GetArtworkForProgramIDs(programIDs []string) ([]ArtworkResponse, error) {
	// If user passed more than 500 programIDs, let's help them out by
	// chunking the requests for them.
	// Obviously you can disable this behavior by passing less than 500 IDs.
	if len(programIDs) > 500 {
		allResponses := make([]ArtworkResponse, 0)
		for _, chunk := range chunkStringSlice(programIDs, 500) {
			resp, err := c.GetArtworkForProgramIDs(chunk)
			if err != nil {
				return nil, err
			}
			allResponses = append(allResponses, resp...)
		}
		return allResponses, nil
	}

	url := fmt.Sprint(c.BaseURL, APIVersion, "/metadata/programs")

	js, jsErr := json.Marshal(programIDs)
	if jsErr != nil {
		return nil, jsErr
	}

	// setup the request
	req, httpErr := http.NewRequest("POST", url, bytes.NewBuffer(js))
	if httpErr != nil {
		return nil, httpErr
	}

	_, data, err := c.SendRequest(req, false)
	if err != nil {
		return nil, err
	}

	programArtwork := make([]ArtworkResponse, 0)

	if err = json.Unmarshal(data, &programArtwork); err != nil {
		return nil, err
	}

	return programArtwork, err
}

// GetArtworkForRootID returns artwork for the given programIDs.
func (c *Client) GetArtworkForRootID(rootID string) ([]Artwork, error) {
	url := fmt.Sprint(c.BaseURL, APIVersion, "/metadata/programs/", rootID)

	// setup the request
	req, httpErr := http.NewRequest("GET", url, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	_, data, err := c.SendRequest(req, false)
	if err != nil {
		return nil, err
	}

	programArtwork := make([]Artwork, 0)

	if err = json.Unmarshal(data, &programArtwork); err != nil {
		return nil, err
	}

	return programArtwork, err
}

// GetImage returns an image for the given URI.
func (c *Client) GetImage(imageURI string) ([]byte, error) {
	url := fmt.Sprint(c.BaseURL, APIVersion, "/image/", imageURI)

	if strings.HasPrefix(imageURI, "https://s3.amazonaws.com") {
		url = imageURI
	}

	// setup the request
	req, httpErr := http.NewRequest("GET", url, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	_, data, err := c.SendRequest(req, true)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetCelebrityArtwork returns artwork for the given programIDs.
func (c *Client) GetCelebrityArtwork(celebrityID string) ([]Artwork, error) {
	url := fmt.Sprint(c.BaseURL, APIVersion, "/metadata/celebrity/", celebrityID)

	// setup the request
	req, httpErr := http.NewRequest("GET", url, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	_, data, err := c.SendRequest(req, false)
	if err != nil {
		return nil, err
	}

	programArtwork := make([]Artwork, 0)

	if err = json.Unmarshal(data, &programArtwork); err != nil {
		return nil, err
	}

	return programArtwork, err
}

// GetImageURL will return a fully formed image URL for the piece given. If it is already fully formed the input will be returned.
func (c *Client) GetImageURL(imageURI string) string {
	if strings.HasPrefix(imageURI, "https://s3.amazonaws.com") {
		return imageURI
	}
	return fmt.Sprint(c.BaseURL, APIVersion, "/image/", imageURI)
}
