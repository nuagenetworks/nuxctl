package nuagex

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	yaml "gopkg.in/yaml.v2"
)

// Lab defines a NuageX environment
type Lab struct {
	ID         string    `yaml:"_id" json:"_id"`
	Name       string    `yaml:"name" json:"name"`
	Reason     string    `yaml:"reason,omitempty" json:"reason"`
	Expires    time.Time `yaml:"expires" json:"expires"`
	Template   string    `yaml:"template" json:"template"`
	SSHKeys    []SSHKey  `yaml:"sshKeys" json:"sshKeys"`
	Services   []Service `yaml:"services" json:"services"`
	Networks   []Network `yaml:"networks" json:"networks"`
	Servers    []Server  `yaml:"servers" json:"servers"`
	Status     string
	Password   string
	ExternalIP string `yaml:"externalIP" json:"externalIP"`
}

// LabResponse : NuageX Lab response JSON object mapping
type LabResponse struct {
	ID       string `json:"_id"`
	Name     string
	Password string
	Status   string
}

// Conf loads nuagex lab configuration from a YAML file
func (c *Lab) Conf(fn string) *Lab {
	fmt.Printf("Loading lab configuration from '%s' file\n", fn)
	yamlFile, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Printf("Lab Configuration Load error   #%v ", err)
		os.Exit(1)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

// CreateLab : Create a Lab in NuageX
func CreateLab(u *User, reqb []byte) (LabResponse, *http.Response, error) {
	URL := buildURL("/labs")
	b, r, err := SendHTTPRequest("POST", URL, u.Token, reqb)
	if err != nil {
		return LabResponse{}, r, err
	}
	var lr LabResponse
	json.Unmarshal(b, &lr)
	if r.StatusCode != 200 {
		var eresp ErrorResponse
		json.Unmarshal(b, &eresp)
		log.Fatalf("Failed to create lab. Reason: %s", eresp.Message)
	}
	return lr, r, nil
}

// GetLab retrieves Lab JSON object
func GetLab(u *User, id string) (Lab, *http.Response, error) {
	URL := buildURL(fmt.Sprintf("/labs/%v?expand=true", id))
	b, r, err := SendHTTPRequest("GET", URL, u.Token, nil)

	if err != nil {
		log.Fatal(err)
	}
	if r.StatusCode != 200 {
		var eresp ErrorResponse
		json.Unmarshal(b, &eresp)
		log.Fatalf("Failed to dump the lab. Reason: %s", eresp.Message)
	}
	var result Lab
	json.Unmarshal(b, &result)
	return result, r, nil
}

// GetLabs retrives Lab JSON objects
func GetLabs(u *User) ([]*Lab, error) {
	URL := buildURL(fmt.Sprintf("/labs?user=%v", u.UserID))
	b, _, err := SendHTTPRequest("GET", URL, u.Token, nil)
	if err != nil {
		log.Fatal(err)
	}
	var l []*Lab
	json.Unmarshal(b, &l)
	return l, nil
}
