package nuagex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// User represents NuageX user
type User struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Token    string
	UserID   string
}

// LoadCredentials : load user credentials from YAML file or env variables
func (u *User) LoadCredentials(fn string) *User {
	// if credentials file exists
	if _, err := os.Stat(fn); err == nil {
		fmt.Printf("Loading user credentials from '%s' file...\n", fn)
		yamlFile, err := ioutil.ReadFile(fn)
		if err != nil {
			log.Printf("LoadCredentials error   #%v ", err)
			os.Exit(1)
		}
		err = yaml.Unmarshal(yamlFile, u)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
		}

	}

	if _, err := os.Stat(fn); os.IsNotExist(err) {
		fmt.Println("Looking for user credentials in the environment variables 'NUX_USERNAME', 'NUX_PASSWORD'...")
		if os.Getenv("NUX_USERNAME") == "" || os.Getenv("NUX_PASSWORD") == "" {
			log.Fatalf("Failed to find credentials file by the default path '%s' and in the environment variables 'NUX_USERNAME', 'NUX_PASSWORD'", fn)
		}
		u.Username, u.Password = os.Getenv("NUX_USERNAME"), os.Getenv("NUX_PASSWORD")
	}
	return u
}

// LoginResponse : represents the credentials data reported back from API server
type LoginResponse struct {
	Token     string `json:"accessToken"`
	User      string `json:"user"`
	ExpiresIn int    `json:"expiresIn"`
}

// Login logs in a user with the passed in `login User` struct
func (u *User) Login() (*User, error) {
	fmt.Printf("Logging '%s' user in...\n", u.Username)
	body, _ := json.Marshal(u)
	URL := buildURL("/auth/login")
	req, _ := http.NewRequest("POST", URL, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	// TODO: rework to use a single HTTPRequest func
	if UA != "" {
		req.Header.Set("User-Agent", UA)
	}
	client := &http.Client{}

	response, _ := client.Do(req)

	if response.StatusCode != 200 {
		fmt.Printf("Failed to login a user! Aborting...\n")
		os.Exit(1)
	}
	defer response.Body.Close()

	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	// fmt.Printf("User '%s' logged in...\n", u.Username)
	var loginResponse LoginResponse

	json.Unmarshal(body, &loginResponse)

	u.Token = loginResponse.Token
	u.UserID = loginResponse.User

	return u, nil
}
