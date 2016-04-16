// Package GoSchedulesDirect provides structs and functions to interact with
// the SchedulesDirect JSON listings service in Go.
package GoSchedulesDirect

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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
	UserAgent      = "TBD"
    ActiveClient *Client
)


// A TokenResponse stores the SD json response message for token request.
type TokenResponse struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	ServerID string `json:"serverID"`
	Token    string `json:"token"`
}

// A VersionResponse stores the SD json response message for a version request.
type VersionResponse struct {
	Response string    `json:"response"`
	Code     int       `json:"code"`
	Client   string    `json:"client"`
	Version  string    `json:"version,omitempty"`
	ServerID string    `json:"serverID"`
	DateTime time.Time `json:"datetime"`
}

// An AddLineupResponse stores the SD json message returned after attempting
// to add a lineup.
type AddLineupResponse struct {
	Response         string    `json:"response"`
	Code             int       `json:"code"`
	ServerID         string    `json:"serverID"`
	Message          string    `json:"message"`
	ChangesRemaining int       `json:"changesRemaining"`
	DateTime         time.Time `json:"datetime"`
}

// A LineupResponse stores the SD json message returned after requesting
// to list subscribed lineups.
type LineupResponse struct {
	ServerID string    `json:"serverID"`
	DateTime time.Time `json:"datetime"`
	Lineups  []Lineup  `json:"lineups"`
}

// A StatusResponse stores the SD json message returned after requesting system
// status.  SystemStatus[0].Status should be "Online" before proceeding.
type StatusResponse struct {
	Account        AccountInfo `json:"account"`
	Lineups        []Lineup    `json:"lineups"`
	LastDataUpdate string      `json:"lastDataUpdate"`
	Notifications  []string    `json:"notifications"`
	SystemStatus   []Status    `json:"systemStatus"`
	ServerID       string      `json:"serverID"`
	Code           int         `json:"code"`
}

// A StatusError struct stores the error response to a status request.
type StatusError struct {
    Response string `json:"response"`
    Code     int    `json:"code"`
    ServerID string `json:"serverID"`
    Message  string `json:"message"`
    Datetime string `json:"datetime"`
    Token    string `json:"token"`
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

// A ChannelResponse stores the channel response for a lineup
type ChannelResponse struct {
	Map      []ChannelMap        `json:"map"`
	Stations []Station           `json:"stations"`
	Metadata ChannelResponseMeta `json:"metadata"`
}

// A ChannelResponseMeta stores the metadata field associated with a channel response
type ChannelResponseMeta struct {
	Lineup    string    `json:"lineup"`
	Modified  time.Time `json:"modified"`
	Transport string    `json:"transport"`
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

// A BroadcasterInfo stores the information about a broadcaster.
type BroadcasterInfo struct {
    City                string      `json:"city"`
    State               string      `json:"state"`
    Postalcode          string      `json:"postalcode"`
    Country             string      `json:"country"`
}

// A StationLogo stores the information to locate a station logo
type StationLogo struct {
	URL    string `json:"URL"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
	Md5    string `json:"md5"`
}

// A ChannelMap stores the station id to channel mapping
type ChannelMap struct {
	StationID string `json:"stationID"`
	Channel   string `json:"channel,omitempty"`
	UhfVhf    int    `json:"uhfVhf,omitempty"`
	AtscMajor int    `json:"atscMajor,omitempty"`
	AtscMinor int    `json:"atscMinor,omitempty"`
}

// A Schedule stores the program information for a given stationID
type Schedule struct {
	StationID string       `json:"stationID"`
	MetaData  ScheduleMeta `json:"metadata"`
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

// A Program stores the information to describing a single television program.
type Program struct {
	ProgramID           string          `json:"programID,omitempty"`
	AirDateTime         time.Time       `json:"airDateTime,omitempty"`
	Md5                 string          `json:"md5,omitempty"`
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
	Ratings             []ProgramRating `json:"ratings, omitempty"`
	ProgramPart         Part            `json:"multipart, omitempty"`
	VideoProperties     []string        `json:"videoProperties,omitempty"`
}

// A SyndicationType stores syndication information for a program
type SyndicationType struct {
	Source string `json:"source"`
	Type   string `json:"type"`
}

// A ProgramRating stores ratings board information for a program
type ProgramRating struct {
	Body string `json:"body"`
	Code string `json:"code"`
}

// A ProgramMetaItem stores meta information for a program
type ProgramMetaItem struct {
	Season  int `json:"season"`
	Episode int `json:"episode,omitmepty"`
}

// A ProgramInfo type stores program information for a program
type ProgramInfo struct {
	ProgramID string `json:"programID"`
	Titles    []struct {
		Title120 string `json:"title120"`
	} `json:"titles"`

	EventDetails    Details                    `json:"eventDetails"`
	Descriptions    ProgramDescriptions        `json:"descriptions"`
	OriginalAirDate string                     `json:"originalAirDate"`
	EpisodeTitle150 string                     `json:"episodeTitle150"`
	Metadata        map[string]ProgramMetaItem `json:"metadata"`
	Movie           Movie                      `json:"movie,omitempty"`
	Cast            []Person                   `json:"cast"`
	Crew            []Person                   `json:"crew"`
	ShowType        string                     `json:"showType"`
	HasImageArtWork bool                       `json:"hasImageArtwork"`
	Md5             string                     `json:"md5"`
}

// A Movie type stores information about a movie
type Movie struct {
	Duration      int    `json:"duration"`
	Year          string `json:"year"`
	QualityRating []struct {
		Increment   string `json:"increment"`
		MaxRating   string `json:"maxRating"`
		MinRating   string `json:"minRating"`
		Rating      string `json:"rating"`
		RatingsBody string `json:"ratingsBody"`
	} `json:"qualityRating"`
}

// Person stores information for an acting credit.
type Person struct {
	PersonID      string `json:"personId,omitmepty"`
	NameID        string `json:"nameId,omitempty"`
	Name          string `json:"name"`
	Role          string `json:"role"`
	CharacterName string `json:"characterName,omitempty"`
	BillingOrder  string `json:"billingOrder"`
}

// ProgramInfoError stores the error response for a program request
type ProgramInfoError struct {
	Response string    `json:"reponse"`
	Code     int       `json:"code"`
	ServerID string    `json:"serverID"`
	Message  string    `json:"message"`
	DateTime time.Time `json:"datetime"`
}

// Details stores some details hehehe
type Details struct {
	Subtype string
}

// ProgramDescriptions stores the descriptive summaries for a program
type ProgramDescriptions struct {
	Description100  []Description `json:"description100,omitempty"`
	Description1000 []Description `json:"description1000.omitempty"`
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

// LastmodifiedRequest stores the information needed to make a last modified request.
type LastmodifiedRequest struct {
	StationID string `json:"stationID"`
	Days      int    `json:"days"`
}

// Client type
type Client struct {
	//Our HTTP client to communicate with SD
	//client *http.Client

	//The Base URL for SD requests
	BaseURL *url.URL
  
    // HTTP
	HTTP *http.Client
	
    //The token 
    Token string
    
	//User agent string
	UserAgent string
}

// NewClient returns a new SD API client.  Uses http.DefaultClient if no
// client is provided.
// TODO Add userAgent string once determined
func NewClient(username string, password string) *Client {
	baseURL, _ := url.Parse(DefaultBaseURL)
	c := &Client{HTTP: &http.Client{}, BaseURL: baseURL}
    token, _ := c.GetToken(username, password)
    c.Token = token
    ActiveClient = c
	return c
}

// encryptPassword returns the sha1 hex encoding of the string argument
func encryptPassword(password string) string {
	encoded := sha1.New()
	encoded.Write([]byte(password))
	return hex.EncodeToString(encoded.Sum(nil))
}

// GetToken returns a session token if the supplied username/password
// successfully authenticate.
func (c Client) GetToken(username string, password string) (string, error) {
	//The SchedulesDirect token url
	url := fmt.Sprint(DefaultBaseURL, APIVersion, "/token")

	//encrypt the password
	sha1hexPW := encryptPassword(password)

	//TODO: Evaluate performance of this string concatenation, not that this
	//should run often.
	var jsonStr = []byte(
		`{"username":"` + username +
			`", "password":"` + sha1hexPW + `"}`)

	//setup the request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	//perform the POST
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close() //resp.Body.Close() will run when we're finished.

	//create a TokenResponse struct, return if err
	r := new(TokenResponse)

	//decode the response body into the new TokenResponse struct
	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		return "", err
	}

	//Print some debugging output
	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body:", string(body))

	//return the token string
	return r.Token, nil
}

// GetStatus returns a StatusResponse for this account.
func (c Client) GetStatus() (*StatusResponse, error) {
	url := fmt.Sprint(DefaultBaseURL, APIVersion, "/status")
	fmt.Println("URL:>", url)
	s := new(StatusResponse)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("token", c.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return s, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatal(resp.Status)
		return s, err
	}
	defer resp.Body.Close() //resp.Body.Close() will run when we're finished.

	//Copy the body to Stdout
	//_, err = io.Copy(os.Stdout, resp.Body)

	err = json.NewDecoder(resp.Body).Decode(s)
	if err != nil {
		fmt.Println("Error parsing status response")
		log.Fatal(err)
		return s, err
	}
	//fmt.Println("Current Status is: ")
	//fmt.Println(s.SystemStatus[0].Status)
	return s, nil
}

// AddLineup adds the given lineup uri to the users SchedulesDirect account.
func (c Client) AddLineup(lineupURI string) error {
	//url := "https://json.schedulesdirect.org" + lineupURI
	url := fmt.Sprint(DefaultBaseURL, lineupURI)
	fmt.Println("URL:>", url)

	req, err := http.NewRequest("PUT", url, nil)
	req.Header.Set("token", c.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)

		return err
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatal(resp.Status)
		return err
	}
	defer resp.Body.Close() //resp.Body.Close() will run when we're finished.

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Add Lineup Response:", string(body))
	return nil
}

// DelLineup deletes the given lineup uri from the users SchedulesDirect account.
func (c Client) DelLineup(lineupURI string) error {
	//url := "https://json.schedulesdirect.org" + lineupURI
	url := fmt.Sprint(DefaultBaseURL, lineupURI)
	fmt.Println("URL:>", url)

	req, err := http.NewRequest("DELETE", url, nil)
	req.Header.Set("token", c.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return err
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatal(resp.Status)
		return err
	}
	defer resp.Body.Close() //resp.Body.Close() will run when we're finished.

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Delete Lineup Response:	Premiere //- Should only be found in Miniseries and Movie program types.", string(body))
	return nil
}

// GetHeadends returns the map of headends for the given postal code.
func (c Client) GetHeadends(postalCode string) ([]Headend, error) {
	url := fmt.Sprint(DefaultBaseURL, APIVersion,
		"/headends?country=USA&postalcode=", postalCode)
	fmt.Println("URL:>", url)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("token", c.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatal(resp.Status)
		return nil, err
	}
	defer resp.Body.Close() //resp.Body.Close() will run when we're finished.

	//make the slice of headends
	h := []Headend{}
    
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("PostalResponse Body:", string(body))

	//decode the body
	err = json.NewDecoder(resp.Body).Decode(&h)
	if err != nil {
		fmt.Println("Error parsing headend responseline")
		log.Fatal(err)
		return nil, err
	}
	return h, nil
}

// GetChannels returns the channels in a given lineup
func (c Client) GetChannels(lineupURI string) (*ChannelResponse, error) {
	url := fmt.Sprint(DefaultBaseURL, lineupURI)
	fmt.Println("URL:>", url)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("token", c.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatal(resp.Status)
		return nil, err
	}
	defer resp.Body.Close() //resp.Body.Close() will run when we're finished.

	//make the map
	h := new(ChannelResponse)

    //debug code	
    //body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))

	//decode the body into the map
	err = json.NewDecoder(resp.Body).Decode(&h)
	if err != nil {
		fmt.Println("Error parsing channel response line")
		log.Fatal(err)
		return nil, err
	}

	return h, nil
}

// GetSchedules returns the set of schedules requested.  As a whole the response is not valid json but each individual line is valid.
func (c Client) GetSchedules(stationIds []string, dates []string) ([]Schedule, error) {
	url := fmt.Sprint(DefaultBaseURL, APIVersion, "/schedules")
	fmt.Println("URL:>", url)

	//buffer to store the json request
	var buffer bytes.Buffer

	//creating the request
	buffer.WriteString("[")
	for index, station := range stationIds {
		//fmt.Println(station)
		buffer.WriteString(`{"stationID":"`+ station + `","date":[`)
        	for index2, date := range dates {
		        buffer.WriteString(`"`+date+`"`)
		        if index2 != len(dates)-1 {
			        buffer.WriteString(",")
		        } else {
                    buffer.WriteString("]")
                }
            }
		if index != len(stationIds)-1 {
			buffer.WriteString("},")
		} else {
            buffer.WriteString("}")
        }
	}
	buffer.WriteString("]")

	//setup the request
	req, err := http.NewRequest("POST", url, &buffer)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Encoding", "deflate,gzip")
	req.Header.Set("token", c.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatal(resp.Status)
		return nil, err
	}
	defer resp.Body.Close() //resp.Body.Close() will run when we're finished.
    
    //decode the response
	var h []Schedule
             
   //debug code	
    //body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
    
	//decode the body
	err = json.NewDecoder(resp.Body).Decode(&h)
	if err != nil {
		fmt.Println("Error parsing schedules response")
		log.Fatal(err)
		return nil, err
	}

	return h, nil
}

// GetProgramInfo returns the set of program details for the given set of programs
func (c Client) GetProgramInfo(programIDs []string) ([]ProgramInfo, error) {
	url := fmt.Sprint(DefaultBaseURL, APIVersion, "/programs")
	fmt.Println("URL:>", url)

	//buffer to store the json request
	var buffer bytes.Buffer

	//creating the request
	buffer.WriteString("[")
	for index, program := range programIDs {
		//fmt.Println(station)
		buffer.WriteString(`"`)
		buffer.WriteString(program)
		buffer.WriteString(`"`)
		if index != len(programIDs)-1 {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString("]")

	//setup the request
	req, err := http.NewRequest("POST", url, &buffer)
	//req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Encoding", "deflate")
	req.Header.Set("token", c.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatal(resp.Status)
		return nil, err
	}
	defer resp.Body.Close() //resp.Body.Close() will run when we're finished.

	//Copy the body to Stdout
	//_, err = io.Copy(os.Stdout, resp.Body)

	//create the schedules slice
	var allPrograms []ProgramInfo

	//readbuffer := bytes.NewBuffer(resp.Body)
	reader := bufio.NewReader(resp.Body)

	//we need to increase the default reader size to get this in one shot
	bufio.NewReaderSize(reader, 65536)
	// there are a few possible loop termination
	// conditions, so just start with an infinite loop.
	for {
		//ReadString because Schedules Direct puts each schedule on it's own line
		//each line is a complete json object but not the whole response.
		line, err := reader.ReadString('\n')

		//debug
		fmt.Println(line)

		// loop termination condition 1:  EOF.
		// this is the normal loop termination condition.
		if err == io.EOF {
			break
		}

		// loop termination condition 2: some other error.
		// Errors happen, so check for them and do something with them.
		if err != nil {
			log.Fatal(err)
		}

		//create a Schedule variable
		var s ProgramInfo

		//decode the scanner bytes into the schedule
		errUnmarshal := json.Unmarshal([]byte(line), &s)
		if errUnmarshal != nil {
			log.Printf("error unmarshaling program: %s\n", errUnmarshal)
		} else {
			allPrograms = append(allPrograms, s)
		}
	}

	return allPrograms, err
}

// GetLastModified returns 
func (c Client) GetLastModified(theRequest []LastmodifiedRequest) {
	url := fmt.Sprint(DefaultBaseURL, APIVersion, "/schedules/md5")
	fmt.Println("URL:>", url)

}

// GetLineups returns a LineupResponse which contains all the lineups subscribed
// to by this account.
func (c Client) GetLineups() (*LineupResponse, error) {
	url := fmt.Sprint(DefaultBaseURL, APIVersion, "/lineups")
	fmt.Println("URL:>", url)
	s := new(LineupResponse)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("token", c.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return s, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatal(resp.Status)
		return s, err
	}
	defer resp.Body.Close() //resp.Body.Close() will run when we're finished.

	//Copy the body to Stdout
	//_, err = io.Copy(os.Stdout, resp.Body)

	err = json.NewDecoder(resp.Body).Decode(s)
	if err != nil {
		fmt.Println("Error parsing status response")
		log.Fatal(err)
		return s, err
	}
	return s, nil
}

