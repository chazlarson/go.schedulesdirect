package schedulesdirect

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Service describes an available service such as countries.
type Service struct {
	Description string `json:"description"`
	Type        string `json:"type"`
	URI         string `json:"uri"`
}

// Country describes a country that Schedules Direct supports.
type Country struct {
	FullName          string `json:"fullName"`
	PostalCode        string `json:"postalCode"`
	PostalCodeExample string `json:"postalCodeExample"`
	ShortName         string `json:"shortName"`
	OnePostalCode     bool   `json:"onePostalCode"`
}

// AvailableDVBS is a single satellite available via DVB-S.
type AvailableDVBS struct {
	Lineup string `json:"lineup"`
}

// GetAvailableServices returns the available services.
func (c *Client) GetAvailableServices() ([]Service, error) {
	url := fmt.Sprint(c.BaseURL, APIVersion, "/available")

	req, httpErr := http.NewRequest("GET", url, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	_, data, err := c.SendRequest(req, false)
	if err != nil {
		return nil, err
	}

	r := make([]Service, 0)

	err = json.Unmarshal(data, &r)
	return r, err
}

// GetAvailableCountries returns the list of countries, grouped by region, supported by Schedules Direct.
func (c *Client) GetAvailableCountries() (map[string][]Country, error) {
	url := fmt.Sprint(c.BaseURL, APIVersion, "/available/countries")

	req, httpErr := http.NewRequest("GET", url, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	_, data, err := c.SendRequest(req, false)
	if err != nil {
		return nil, err
	}

	r := make(map[string][]Country)

	err = json.Unmarshal(data, &r)

	return r, err
}

// GetAvailableLanguages returns the list of language digraphs and their language names.
func (c *Client) GetAvailableLanguages() (map[string]string, error) {
	url := fmt.Sprint(c.BaseURL, APIVersion, "/available/languages")

	req, httpErr := http.NewRequest("GET", url, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	_, data, err := c.SendRequest(req, false)
	if err != nil {
		return nil, err
	}

	r := make(map[string]string)

	err = json.Unmarshal(data, &r)

	return r, err
}

// GetAvailableDVBS returns the list of satellites which are available.
func (c *Client) GetAvailableDVBS() ([]AvailableDVBS, error) {
	url := fmt.Sprint(c.BaseURL, APIVersion, "/available/dvb-s")

	req, httpErr := http.NewRequest("GET", url, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	_, data, err := c.SendRequest(req, false)
	if err != nil {
		return nil, err
	}

	r := make([]AvailableDVBS, 0)

	err = json.Unmarshal(data, &r)

	return r, err
}

// GetAvailableTransmitters returns the list of freeview transmitters in a country for the given countryCode.
// Country options: GBR.
func (c *Client) GetAvailableTransmitters(countryCode string) (map[string]string, error) {
	url := fmt.Sprint(c.BaseURL, APIVersion, "/available/transmitters/", countryCode)

	req, httpErr := http.NewRequest("GET", url, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	_, data, err := c.SendRequest(req, false)
	if err != nil {
		return nil, err
	}

	r := make(map[string]string)

	err = json.Unmarshal(data, &r)

	return r, err
}
