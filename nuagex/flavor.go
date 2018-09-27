package nuagex

import (
	"encoding/json"
	"log"
)

// Flavor defines a NuageX flavor for image
type Flavor struct {
	ID     string `yaml:"_id,omitempty" json:"_id"`
	Name   string `yaml:"name" json:"name"`
	CPU    int    `yaml:"cpu" json:"cpu"`
	Memory int    `yaml:"memory" json:"memory"`
	Disk   int    `yaml:"disk" json:"disk"`
}

// GetFlavors retrives Flavor JSON objects
func GetFlavors(u *User) ([]*Flavor, error) {
	URL := buildURL("/flavors?limit=100")
	b, _, err := SendHTTPRequest("GET", URL, u.Token, nil)
	if err != nil {
		log.Fatal(err)
	}
	var f []*Flavor
	json.Unmarshal(b, &f)
	return f, nil
}
