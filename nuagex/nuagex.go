package nuagex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const apiURL string = "https://experience.nuagenetworks.net/api"

// UA is a string to receive custom User-Agent for HTTP requests
var UA string

// FlavorResponse : Image Flavor response JSON object mapping
// NOT IMPLEMENTED YET IN CLI
type FlavorResponse struct {
	ID                   string        `json:"_id"`
	Name                 string        `json:"name"`
	CPU                  int           `json:"cpu"`
	Memory               int           `json:"memory"`
	V                    int           `json:"__v"`
	AllowedGroups        []string      `json:"allowedGroups"`
	AllowedOrganizations []interface{} `json:"allowedOrganizations"`
	Disk                 int           `json:"disk"`
	Created              time.Time     `json:"created"`
}

// ErrorResponse handles NuageX API responce in case of a non 200 OK operation
type ErrorResponse struct {
	Message string
}

func buildURL(path string) string {
	return apiURL + path
}

// SendHTTPRequest : Request for a given method and url along with body parameters.
func SendHTTPRequest(method, url, token string, body []byte) ([]byte, *http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	if UA != "" {
		req.Header.Set("User-Agent", UA)
	}

	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	return body, response, nil
}

// GetAllFlavors : Get all image flavors
func GetAllFlavors(token string) ([]FlavorResponse, error) {
	URL := buildURL("/flavors")
	b, _, err := SendHTTPRequest("GET", URL, token, nil)

	if err != nil {
		return []FlavorResponse{}, err
	}

	var result []FlavorResponse

	json.Unmarshal(b, &result)

	return result, nil
}
