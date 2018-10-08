package cmd

import (
	"fmt"
	"os"

	"github.com/nuagenetworks/nuxctl/nuagex"

	"github.com/spf13/cobra"
)

// VERSION is set in main.go and tells the nuxctl version
var VERSION string

var emptyTemplateID = "5980ee745a38da00012d158d"
var nuxReason = "nuxctl"

var user nuagex.User

var lab nuagex.Lab

// labID is the ID of a NuageX Lab (which is also the hostname name)
var labID string

var rootCmd = &cobra.Command{
	Use:   "nuxctl",
	Short: "nuxctl is a CLI client for NuageX lab deployment",
	Long:  `nuxctl is a command line client to deploy labs which configuration is expressed in YAML files on the NuageX platform.`,
}

// Execute launches the root command
func Execute(ver string) {
	VERSION = ver
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
