package nuagex

import (
	"encoding/json"
	"log"
)

// Image defines a NuageX image
type Image struct {
	ID      string `yaml:"_id,omitempty" json:"_id"`
	Name    string `yaml:"name" json:"name"`
	MinDisk int    `yaml:"minDisk" json:"minDisk"`
}

// GetImages retrives Image JSON objects
func GetImages(u *User) ([]*Image, error) {
	URL := buildURL("/images?limit=100")
	b, _, err := SendHTTPRequest("GET", URL, u.Token, nil)
	if err != nil {
		log.Fatal(err)
	}
	var i []*Image
	json.Unmarshal(b, &i)
	return i, nil
}
