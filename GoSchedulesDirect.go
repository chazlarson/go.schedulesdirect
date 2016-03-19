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
	//"os"
	"strconv"
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
const (
	apiVersion     = "20141201"
	defaultBaseURL = "https://json.schedulesdirect.org/"
	userAgent      = "TBD"
)

//Client type
type Client struct {
	//Our HTTP client to communicate with SD
	client *http.Client

	//The Base URL for SD requests
	BaseURL *url.URL

	//User agent string
	UserAgent string
}

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
	ServerId       string      `json:"serverID"`
	code           int         `json:"code"`
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

// A ChannelReponse stores the channel response for a lineup
type ChannelResponse struct {
	Map      []ChannelMap        `json:"map"`
	Stations []Station           `json:"stations"`
	Metadata ChannelResponseMeta `json:"metadata"`
}

// A ChannelResponseMeta stores the metadata field associated with a channel response
type ChannelResponseMeta struct {
	Lineup    string    `json:"lineup"`
	Modified  time.Time `json:"modified"`
	Transport string    `json"transport"`
}

// A Station stores the SD json that describes a station.
type Station struct {
	Callsign            string      `json:"callsign"`
	IsCommercialFree    bool        `json:"isCommercialFree"`
	Name                string      `json:"name"`
	BroadcastLanguage   string      `json:"broadcastLanguage"`
	DescriptionLanguage string      `json:"descriptionLanguage "`
	Logo                StationLogo `json:"logo"`
	StationId           string      `json:"stationID"`
}

// A StationLogo stores the information to locate a station logo
type StationLogo struct {
	Url    string `json:"URL"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
	Md5    string `json:"md5"`
}

// A ChannelMap stores the station id to channel mapping
type ChannelMap struct {
	StationId string `json:"stationID"`
	Channel   int    `json:"channel,omitempty"`
	UhfVhf    int    `json:"uhfVhf,omitempty"`
	AtscMajor int    `json:"atscMajor,omitempty"`
	AtscMinor int    `json:"atscMinor,omitempty"`
}

type ScheduleResponse struct {
	Schedules []Schedule
}

// A Schedule stores the program information for a given stationID
type Schedule struct {
	StationId string       `json:"stationID"`
	MetaData  ScheduleMeta `json:"metadata"`
	Programs  []Program    `json:"programs"`
}

type ScheduleMeta struct {
	Modified  string `json:"modified"`
	MD5       string `json:"md5"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Days      int    `json:"days"`
}

// A Program stores the information to describing a single television program.
type Program struct {
	ProgramId           string          `json:"programID"`
	AirDateTime         time.Time       `json:"airDateTime"`
	Md5                 string          `json:"md5"`
	Duration            int             `json:"duration"`
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
	Syndication         SyndicationType `json:"syndication"`
	Ratings             []ProgramRating `json:"ratings, omitempty"`
	ProgramPart         Part            `json:"multipart, omitempty"`
	VideoProperties     []string        `json:"videoProperties,omitempty"`
}
type SyndicationType struct {
	Source string `json:"source"`
	Type   string `json:"type"`
}

type ProgramRating struct {
	Body string `json:"body"`
	Code string `json:"code"`
}

type ProgramMetaItem struct {
	Season  int `json:"season"`
	Episode int `json:"episode,omitmepty"`
}

type ProgramInfo struct {
	ProgramId string `json:"programID"`
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

type Person struct {
	PersonID      string `json:"personId,omitmepty"`
	NameID        string `json:"nameId,omitempty"`
	Name          string `json:"name"`
	Role          string `json:"role"`
	CharacterName string `json:"characterName,omitempty"`
	BillingOrder  string `json:"billingOrder"`
}

type ProgramInfoError struct {
	Response string    `json:"reponse"`
	Code     int       `json:"code"`
	ServerId string    `json:"serverID"`
	Message  string    `json:"message"`
	DateTime time.Time `json:"datetime"`
}

type Details struct {
	Subtype string
}

type ProgramDescriptions struct {
	Description100  []Description `json:"description100,omitempty"`
	Description1000 []Description `json:"description1000.omitempty"`
}

type Description struct {
	DescriptionLanguage string `json:"descriptionLanguage"`
	Description         string `json:"description"`
}
type Part struct {
	PartNumber int `json:"partNumber"`
	TotalParts int `json:"totalParts"`
}

type LastmodifiedRequest struct {
	StationId string `json:"stationID"`
	Days      int    `json:"days"`
}

//NewClient returns a new SD API client.  Uses http.DefaultClient if no
//client is provided.
//TODO Add userAgent string once determined
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{client: httpClient, BaseURL: baseURL}
	return c

}

// AddLineup adds the given lineup uri to the users SchedulesDirect account.
func AddLineup(token string, lineupURI string) error {
	//url := "https://json.schedulesdirect.org" + lineupURI
	url := fmt.Sprint("https://json.schedulesdirect.org/", lineupURI)
	fmt.Println("URL:>", url)

	req, err := http.NewRequest("PUT", url, nil)
	req.Header.Set("token", token)

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
func DelLineup(token string, lineupURI string) error {
	//url := "https://json.schedulesdirect.org" + lineupURI
	url := fmt.Sprint("https://json.schedulesdirect.org/", lineupURI)
	fmt.Println("URL:>", url)

	req, err := http.NewRequest("DELETE", url, nil)
	req.Header.Set("token", token)

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
func GetHeadends(token string, postalCode string) ([]Headend, error) {
	url := fmt.Sprint("https://json.schedulesdirect.org/", apiVersion,
		"/headends?country=USA&postalcode=", postalCode)
	fmt.Println("URL:>", url)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("token", token)

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
	//h := make(map[string]Headend)
	h := []Headend{}

	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("PostalResponse Body:", string(body))

	//decode the body into the map
	err = json.NewDecoder(resp.Body).Decode(&h)
	if err != nil {
		fmt.Println("Error parsing headend responseline")
		log.Fatal(err)
		return nil, err
	}
	return h, nil
}

//GetChannels returns the channels in a given lineup
func GetChannels(token string, lineupURI string) (*ChannelResponse, error) {
	url := fmt.Sprint("https://json.schedulesdirect.org", lineupURI)
	fmt.Println("URL:>", url)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("token", token)

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

	//decode the body into the map
	err = json.NewDecoder(resp.Body).Decode(&h)
	if err != nil {
		fmt.Println("Error parsing channel response line")
		log.Fatal(err)
		return nil, err
	}
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(body)
	return h, nil
}

// GetSchedule returns the schedule requested.
func GetSchedule(token string, stationId string, days int) (*Schedule, error) {
	url := fmt.Sprint("https://json.schedulesdirect.org/", apiVersion, "/schedules")
	fmt.Println("URL:>", url)

	var jsonStr = []byte(`[{"stationID":"` + stationId + `", "days":` + strconv.Itoa(days) + `}]`)

	//debug
	//fmt.Println(string(jsonStr))

	//setup the request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Accept-Encoding", "deflate,gzip")
	req.Header.Set("token", token)

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
	h := new(Schedule)

	//decode the body into the map
	err = json.NewDecoder(resp.Body).Decode(&h)
	if err != nil {
		fmt.Println("Error parsing channel response line")
		log.Fatal(err)
		return nil, err
	}

	return h, nil
}

// GetSchedules returns the set of schedules requested.  As a whole the response is not valid json,
// but each individual line is valid.
func GetSchedules(token string, stationIds []string, days int) ([]Schedule, error) {
	url := fmt.Sprint("https://json.schedulesdirect.org/", apiVersion, "/schedules")
	fmt.Println("URL:>", url)

	//buffer to store the json request
	var buffer bytes.Buffer

	//creating the request
	buffer.WriteString("[")
	for index, station := range stationIds {
		//fmt.Println(station)
		buffer.WriteString(fmt.Sprint(`{"stationID":"`, station, `","days":`, strconv.Itoa(days), `}`))
		if index != len(stationIds)-1 {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString("]")

	//setup the request
	req, err := http.NewRequest("POST", url, &buffer)
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Accept-Encoding", "deflate,gzip")
	req.Header.Set("token", token)

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

	//create the schedules slice
	var allSchedules []Schedule

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
		//fmt.Println(line)

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
		var s Schedule

		//decode the scanner bytes into the schedule
		errUnmarshal := json.Unmarshal([]byte(line), &s)
		if errUnmarshal != nil {
			log.Printf("error unmarshaling program: %s\n", errUnmarshal)
			var e ProgramInfoError
			getError := json.Unmarshal([]byte(line), &e)
			if getError != nil {
				fmt.Println("error")
			}

			fmt.Println(e.Code)
		} else {
			allSchedules = append(allSchedules, s)
		}
	}

	return allSchedules, err
}

// GetProgramInfo returns the set of program details for the given set of programs
func GetProgramInfo(token string, programIDs []string) ([]ProgramInfo, error) {
	url := fmt.Sprint("https://json.schedulesdirect.org/", apiVersion, "/programs")
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
	req.Header.Set("token", token)

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

func GetLastModified(token string, theRequest []LastmodifiedRequest) {
	url := fmt.Sprint("https://json.schedulesdirect.org/", apiVersion, "/schedules/md5")
	fmt.Println("URL:>", url)

}

// GetLineups returns a LineupResponse which contains all the lineups subscribed
// to by this account.
func GetLineups(token string) (*LineupResponse, error) {
	url := fmt.Sprint("https://json.schedulesdirect.org/", apiVersion, "/lineups")
	fmt.Println("URL:>", url)
	s := new(LineupResponse)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("token", token)

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

// GetStatus returns a StatusResponse for this account.
func GetStatus(token string) (*StatusResponse, error) {
	url := fmt.Sprint("https://json.schedulesdirect.org/", apiVersion, "/status")
	fmt.Println("URL:>", url)
	s := new(StatusResponse)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("token", token)

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

// GetToken returns a session token if the supplied username/password
// successfully authenticate.
func GetToken(username string, password string) (string, error) {
	//The SchedulesDirect token url
	url := fmt.Sprint("https://json.schedulesdirect.org/", apiVersion, "/token")

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

// encryptPassword returns the sha1 hex enconding of the string argument
func encryptPassword(password string) string {
	encoded := sha1.New()
	encoded.Write([]byte(password))
	return hex.EncodeToString(encoded.Sum(nil))
}
