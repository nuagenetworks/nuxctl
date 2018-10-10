package cmd

import (
	"fmt"
	"log"

	"github.com/nuagenetworks/nuxctl/nuagex"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteLabCmd)

	deleteLabCmd.Flags().StringVarP(&labID, "lab-id", "i", "", "Lab ID")
	deleteLabCmd.MarkFlagRequired("lab-id")

	deleteLabCmd.Flags().StringVarP(&CredFPath, "credentials", "c", "user_creds.yml", "Path to the user credentials file")
}

var deleteLabCmd = &cobra.Command{
	Use:   "delete-lab",
	Short: "Delete existing NuageX lab",
	Long:  `Deletes an existing NuageX lab`,
	Run:   deleteLab,
}

func deleteLab(cmd *cobra.Command, args []string) {
	loginUser(cmd, args)

	fmt.Printf("Deleting NuageX Lab ID %v\n", labID)

	_, _, err := nuagex.DeleteLab(&user, labID)
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Printf("Lab ID '%s' has been scheduled for deletion!\n", labID)
}
