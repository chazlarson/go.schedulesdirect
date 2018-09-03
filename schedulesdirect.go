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
	"time"
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
	// Time we last received a token + 24 hours.
	TokenExpiresAt time.Time

	// The User-Agent to send on every request.
	UserAgent string

	// We store username and password in the client in case we need to attempt a token refresh.
	username       string
	password       string
	failedRequests int
}

// NewClient returns a new Schedules Direct API client. Uses http.DefaultClient if no http.Client is set.
func NewClient(username string, password string) (*Client, error) {
	c := &Client{
		BaseURL:   DefaultBaseURL,
		HTTP:      http.DefaultClient,
		UserAgent: DefaultUserAgent,
		username:  username,
		password:  password,
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

	response, httpErr := c.HTTP.Do(req)
	if httpErr != nil {
		return "", fmt.Errorf("cannot reach schedules direct service: %s", httpErr)
	}

	buf := &bytes.Buffer{}
	if _, copyErr := io.Copy(buf, response.Body); copyErr != nil {
		return "", fmt.Errorf("error when copying bytes of response to buffer: %s", copyErr)
	}

	if closeErr := response.Body.Close(); closeErr != nil {
		return "", fmt.Errorf("cannot read response. %v", closeErr)
	}

	// create a TokenResponse struct, return if err
	r := &TokenResponse{}

	// decode the response body into the new TokenResponse struct
	if err := json.Unmarshal(buf.Bytes(), r); err != nil {
		return "", err
	}

	if r.BaseResponse.Code != ErrOK {
		return "", r.BaseResponse
	}

	c.TokenExpiresAt = r.BaseResponse.DateTime.Add(24 * time.Hour)

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

// SendRequest will send the given http.Request to Schedules Direct.
// Specify if the request requires a token via the needsToken boolean.
func (c Client) SendRequest(request *http.Request, needsToken bool) (*http.Response, []byte, error) {
	if needsToken && c.Token == "" {
		return nil, nil, fmt.Errorf("schedules direct client has not been initialized with a token, stubbornly refusing to make a request")
	}

	// If we've had the token for more than 24 hours we need to refresh it.
	if time.Now().After(c.TokenExpiresAt) && c.failedRequests == 0 {
		c.failedRequests = c.failedRequests + 1
		token, tokenErr := c.GetToken(c.username, c.password)
		if tokenErr != nil {
			return nil, nil, fmt.Errorf("error when attempting to automatically refresh schedules direct token after its expiration: %s", tokenErr)
		}
		c.Token = token
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
		return nil, nil, fmt.Errorf("cannot reach schedules direct service: %s", httpErr)
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
		if baseResp.Code == ErrInvalidUser && c.failedRequests == 0 {
			// We know that at some point the credentials were valid, so let's try running the same request again
			// after we attempt to update the token in case it was expired due to something other than expiration.
			c.failedRequests = c.failedRequests + 1
			token, tokenErr := c.GetToken(c.username, c.password)
			if tokenErr != nil {
				return nil, nil, fmt.Errorf("error when attempting to automatically refresh schedules direct token due to caught INVALID_USER (4003): %s", tokenErr)
			}
			c.Token = token
			return c.SendRequest(request, needsToken)
		} else if baseResp.Code != 0 {
			return nil, nil, baseResp
		}
	}

	if response.StatusCode > 399 {
		return nil, nil, fmt.Errorf("status code was %d, expected 2XX-3XX. received content: %s", response.StatusCode, buf.String())
	}

	c.failedRequests = 0

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
