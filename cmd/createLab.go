package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nuagenetworks/nuxctl/nuagex"

	"github.com/spf13/cobra"
)

// LabFPath is a path to the lab definition file
var LabFPath string

// wait is set when create-lab command should wait till the lab fails or succeeds to deploy
var wait bool

func init() {
	rootCmd.AddCommand(createLabCmd)

	createLabCmd.Flags().StringVarP(&CredFPath, "credentials", "c", "user_creds.yml", "Path to the user credentials file")

	createLabCmd.Flags().StringVarP(&LabFPath, "lab-configuration", "l", "lab.yml", "Path to the Lab configuration file")

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
	lab.Conf(LabFPath)

	lab.Reason = nuxReason // change reason field to nuxctl

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
