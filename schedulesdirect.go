// Package GoSchedulesDirect provides structs and functions to interact with
// the SchedulesDirect JSON listings service in Go.
package GoSchedulesDirect

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	Type     string   `json:"type"`
	Location string   `json:"location"`
	Lineups  []Lineup `json:"lineups"`
}

// A Lineup stores the SD json message containing lineup information.
type Lineup struct {
	Lineup    string `json:lineup,omitempty"`
	Name      string `json:"name,omitempty"`
	ID        string `json:"ID,omitempty"`
	Modified  string `json:"modified,omitempty"`
	URI       string `json:"uri"`
	IsDeleted bool   `json:"isDeleted,omitempty"`
}

// AddLineup adds the given lineup uri to the users SchedulesDirect account.
func AddLineup(token string, lineupURI string) error {
	url := "https://json.schedulesdirect.org" + lineupURI
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
	url := "https://json.schedulesdirect.org" + lineupURI
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
	fmt.Println("Delete Lineup Response:", string(body))
	return nil
}

// GetHeadends returns the map of headends for the given postal code.
func GetHeadends(token string, postalCode string) (map[string]Headend, error) {
	url := "https://json.schedulesdirect.org/20140530/" +
		"headends?country=USA&postalcode=" + postalCode
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
	h := make(map[string]Headend)

	//decode the body into the map
	err = json.NewDecoder(resp.Body).Decode(&h)
	if err != nil {
		fmt.Println("Error parsing headend responsine")
		log.Fatal(err)
		return nil, err
	}
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("PostalResponse Body:", string(body))
	return h, nil
}

/*
func GetLineup(token string, lineup string) {
	url := "https://json.schedulesdirect.org/20140530/lineups/" + lineup

}
*/

// GetLineups returns a LineupResponse which contains all the lineups subscribed
// to by this account.

func GetLineups(token string) (*LineupResponse, error) {
	url := "https://json.schedulesdirect.org/20140530/lineups"
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
	url := "https://json.schedulesdirect.org/20140530/status"
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
	url := "https://json.schedulesdirect.org/20140530/token"

	//encrypt the password
	sha1hexPW := encryptPassword(password)

	//TODO Evaluate performance of this string concatenation, not that this
	//     should run often.
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
