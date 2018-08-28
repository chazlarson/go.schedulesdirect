// Package schedulesdirect provides structs and functions to interact with
// the SchedulesDirect JSON listings service in Go. It is only compatible with
// API version 20141201.
package schedulesdirect

import (
	"bytes"
	"compress/gzip"
	"crypto/sha1" // #nosec
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Some constants for use in the library
const (
	APIVersion       = "20141201"
	DefaultBaseURL   = "https://json.schedulesdirect.org/"
	DefaultUserAgent = "go.schedulesdirect (Go-http-client/1.1)"
)

// Client type
type Client struct {
	// The Base URL for Schedules Direct requests
	BaseURL string

	// Our HTTP client to communicate with Schedules Direct
	HTTP *http.Client

	// The token
	Token string

	// The User-Agent to send on every request.
	UserAgent string
}

// NewClient returns a new Schedules Direct API client. Uses http.DefaultClient if no http.Client is set.
func NewClient(username string, password string) (*Client, error) {
	c := &Client{
		BaseURL:   DefaultBaseURL,
		HTTP:      http.DefaultClient,
		UserAgent: DefaultUserAgent,
	}
	token, tokenErr := c.GetToken(username, password)
	if tokenErr != nil {
		return nil, fmt.Errorf("error getting token from schedules direct: %s", tokenErr)
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
func (c *Client) GetToken(username string, password string) (string, error) {
	url := fmt.Sprint(c.BaseURL, APIVersion, "/token")

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
	_, data, err := c.SendRequest(req, false)
	if err != nil {
		return "", err
	}

	// create a TokenResponse struct, return if err
	r := TokenResponse{}

	// decode the response body into the new TokenResponse struct
	err = json.Unmarshal(data, &r)
	if err != nil {
		return "", err
	}

	// return the token string
	return r.Token, nil
}

// GetStatus returns a StatusResponse for this account.
func (c *Client) GetStatus() (*StatusResponse, error) {
	url := fmt.Sprint(c.BaseURL, APIVersion, "/status")

	req, httpErr := http.NewRequest("GET", url, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	_, data, err := c.SendRequest(req, true)
	if err != nil {
		return nil, err
	}

	s := &StatusResponse{}

	if jsonErr := json.Unmarshal(data, &s); jsonErr != nil {
		return nil, jsonErr
	}
	return s, err
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
func (c *Client) SubmitLineup(hdhrLineupJSON []byte, lineupID string) (*BaseResponse, error) {
	url := fmt.Sprint(c.BaseURL, APIVersion, "/map/lineup/", lineupID)

	req, httpErr := http.NewRequest("POST", url, bytes.NewBuffer(hdhrLineupJSON))
	if httpErr != nil {
		return nil, httpErr
	}

	_, data, err := c.SendRequest(req, true)
	if err != nil {
		return nil, err
	}

	baseR := &BaseResponse{}

	err = json.Unmarshal(data, &baseR)

	return baseR, err
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

// DeleteSystemMessage deletes a system message from the status response.
func (c *Client) DeleteSystemMessage(messageID string) error {
	url := fmt.Sprint(c.BaseURL, "/messages/", messageID)

	req, httpErr := http.NewRequest("DELETE", url, nil)
	if httpErr != nil {
		return httpErr
	}

	_, _, err := c.SendRequest(req, true)

	return err
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

// GetArtworkForProgramIDs returns artwork for the given programIDs.
//
// If more than 500 Program IDs are provided, the client will automatically
// chunk the slice into groups of 500 IDs and return all responses to you.
func (c *Client) GetArtworkForProgramIDs(programIDs []string) ([]ProgramArtworkResponse, error) {
	// If user passed more than 500 programIDs, let's help them out by
	// chunking the requests for them.
	// Obviously you can disable this behavior by passing less than 500 IDs.
	if len(programIDs) > 500 {
		allResponses := make([]ProgramArtworkResponse, len(programIDs))
		for _, chunk := range chunkStringSlice(programIDs, 500) {
			resp, err := c.GetArtworkForProgramIDs(chunk)
			if err != nil {
				return nil, err
			}
			allResponses = append(allResponses, resp...)
		}
		return allResponses, nil
	}

	// Artwork endpoint only wants the leftmost 10 characters of the programID.
	// In case users pass the full 14 character string, let's help them out.
	for idx, programID := range programIDs {
		if len(programID) > 10 {
			programIDs[idx] = programID[:10]
		}
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

	programArtwork := make([]ProgramArtworkResponse, 0)

	if err = json.Unmarshal(data, &programArtwork); err != nil {
		return nil, err
	}

	return programArtwork, err
}

// GetArtworkForRootID returns artwork for the given programIDs.
func (c *Client) GetArtworkForRootID(rootID string) ([]ProgramArtwork, error) {
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

	programArtwork := make([]ProgramArtwork, 0)

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
func (c *Client) GetCelebrityArtwork(celebrityID string) ([]ProgramArtwork, error) {
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

	programArtwork := make([]ProgramArtwork, 0)

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

// SendRequest will send the given http.Request to Schedules Direct.
// Specify if the request requires a token via the needsToken boolean.
func (c *Client) SendRequest(request *http.Request, needsToken bool) (*http.Response, []byte, error) {
	if needsToken && (c == nil || c.Token == "") {
		return nil, nil, fmt.Errorf("schedules direct client has not been initialized with a token, stubbornly refusing to make a request")
	}

	request.Header.Set("User-Agent", c.UserAgent)
	if needsToken {
		request.Header.Set("token", c.Token)
	}

	if request.Method == "POST" {
		request.Header.Set("Content-Type", "application/json")
	}

	response, httpErr := c.HTTP.Do(request)
	if httpErr != nil {
		return nil, nil, fmt.Errorf("cannot reach server schedules direct service: %s", httpErr)
	}

	// This is only for getting programs.
	//
	// From the docs:
	// Your client must send an Accept-Encoding that has "deflate,gzip" in it, even though the response will be gzip'ed.
	// This is due to an implementation bug in 20140530 which will be fixed in 20141201.
	//
	// Not actually fixed yet and Go disables automatic decompression if Accept-Encoding is set, so we are stuck doing the decompression ourselves.
	var reader = response.Body
	if response.Header.Get("Content-Encoding") == "gzip" && !response.Uncompressed {
		readerG, errG := gzip.NewReader(reader)
		if errG == nil {
			reader = readerG
		} else {
			return nil, nil, errG
		}
	}

	buf := &bytes.Buffer{}
	if _, copyErr := io.Copy(buf, reader); copyErr != nil {
		return nil, nil, fmt.Errorf("error when copying bytes of response to buffer: %s", copyErr)
	}

	if closeErr := response.Body.Close(); closeErr != nil {
		return nil, nil, fmt.Errorf("cannot read response. %v", closeErr)
	}

	baseResp := &BaseResponse{}
	if unmarshalErr := json.Unmarshal(buf.Bytes(), baseResp); unmarshalErr == nil {
		if baseResp.Code != 0 {
			return nil, nil, baseResp
		}
	}

	if response.StatusCode > 399 {
		return nil, nil, fmt.Errorf("status code was %d, expected 2XX-3XX. received content: %s", response.StatusCode, buf.String())
	}

	return response, buf.Bytes(), nil
}

// chunkStringSlice will return a slice of slice of strings for the given chunkSize.
func chunkStringSlice(sl []string, chunkSize int) [][]string {
	var divided [][]string

	for i := 0; i < len(sl); i += chunkSize {
		end := i + chunkSize

		if end > len(sl) {
			end = len(sl)
		}

		divided = append(divided, sl[i:end])
	}
	return divided
}
