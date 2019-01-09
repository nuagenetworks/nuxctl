package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/nuagenetworks/nuxctl/nuagex"

	"github.com/spf13/cobra"
)

// LabFPath is a path to the lab definition file
var LabFPath string

// wait is set when create-lab command should wait till the lab fails or succeeds to deploy
var wait bool

// labDuration sets the default value for the number of days that a lab without
// an expiration date will last since its deployment date
// defaults to 14 days
var labDuration = "14d"

func init() {
	rootCmd.AddCommand(createLabCmd)

	createLabCmd.Flags().StringVarP(&CredFPath, "credentials", "c", "user_creds.yml", "Path to the user credentials file")

	createLabCmd.Flags().StringVarP(&LabFPath, "lab-configuration", "l", "lab.yml", "Path to the Lab configuration file")

	createLabCmd.Flags().StringVarP(&labDuration, "duration", "d", labDuration, "Lab duration in the format of: M(onths)w(eeks)d(ays)h(ours). Examples: 5d; 2M4d; 12h")

	createLabCmd.Flags().BoolVarP(&wait, "wait", "w", false, "wait till lab succeeds of fails to deploy")
}

var createLabCmd = &cobra.Command{
	Use:    "create-lab",
	Short:  "Create NuageX lab (environment)",
	Long:   `Create NuageX lab using the lab definition supplied in various formats`,
	PreRun: loginUser,
	Run:    createLab,
}

func createLab(cmd *cobra.Command, args []string) {
	lab.LoadConf(LabFPath)

	lab.Reason = nuxReason // change reason field to nuxctl

	d := parseDuration(labDuration)
	lab.Expires = time.Now().Add(d)

	j, err := json.Marshal(lab)
	if err != nil {
		log.Fatalf("%v", err)
	}
	lr, r, err := nuagex.CreateLab(&user, j)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Lab ID %s has been successfully queued for creation! Request ID %s.\n", lr.ID, r.Header.Get("x-request-id"))

	if wait {
		fmt.Println("Deploying lab...")
		for {
			l, _, err := nuagex.GetLab(&user, lr.ID)
			if err != nil {
				log.Fatal(err)
			}
			if l.Status == "started" {
				fmt.Printf("Lab ID %s has been successfully started!\n", lr.ID)
				printLabs([]*nuagex.Lab{&l})
				return
			} else if l.Status == "errored" {
				fmt.Printf("Lab ID %s has failed to deploy!\n", lr.ID)
				return
			}
			time.Sleep(15 * time.Second)
		}
	}
}

// parseDuration parses the duration that comes as a string with tokens to the time.Duration
// example input strings:
//    1M3w - 1 month, 3 weeks
//    3d - 3 days
//    14d5h - 14 days 5 hours
func parseDuration(s string) time.Duration {
	durationRE := regexp.MustCompile(`(?P<months>\d+M)?(?P<weeks>\d+w)?(?P<days>\d+d)?(?P<hours>\d+h)?`)
	names := durationRE.SubexpNames()
	matches := durationRE.FindStringSubmatch(s)
	if matches[0] == "" {
		log.Fatalf("The duration '%v' was not recognized. Valid tokens are M(onths), d(ays), h(ours)", labDuration)
	}
	namedMatches := map[string]string{}
	for i, n := range matches {
		namedMatches[names[i]] = n
	}
	months := convDuration(namedMatches["months"])
	weeks := convDuration(namedMatches["weeks"])
	days := convDuration(namedMatches["days"])
	hours := convDuration(namedMatches["hours"])
	return time.Duration(months*30*24*int(time.Hour) + weeks*7*24*int(time.Hour) + days*24*int(time.Hour) + hours*int(time.Hour))
}

// convDuration converts the duration string that comes as
// XXd or XXM or XXw to integer XX removing the trailing literal
func convDuration(s string) int {
	if len(s) == 0 {
		return 0
	}
	parsed, err := strconv.Atoi(s[:len(s)-1])
	if err != nil {
		return 0
	}
	return parsed
}
