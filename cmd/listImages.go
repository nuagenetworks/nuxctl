package cmd

import (
	"fmt"
	"log"
	"nuxctl/nuagex"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listImagesCmd)

	listImagesCmd.Flags().StringVarP(&CredFPath, "credentials", "c", "user_creds.yml", "Path to the user credentials file.")
}

var listImagesCmd = &cobra.Command{
	Use:   "list-images",
	Short: "Display NuageX images.",
	Long:  `Outputs to console the list of available NuageX images.`,
	Run:   listImages,
}

func listImages(cmd *cobra.Command, args []string) {
	loginUser(cmd, args)

	fmt.Println("Retrieving NuageX Lab images...")

	images, err := nuagex.GetImages(&user)
	if err != nil {
		log.Fatalf("%v", err)
	}

	printImages(images)
}

type byImageName []*nuagex.Image

func (x byImageName) Len() int           { return len(x) }
func (x byImageName) Less(i, j int) bool { return x[i].Name < x[j].Name }
func (x byImageName) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

func printImages(i []*nuagex.Image) {
	const format = "%v\t%v\n"
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	sort.Sort(byImageName(i))
	fmt.Printf("\n")
	fmt.Fprintf(tw, format, "Name", "Min. Disk (GB)")
	fmt.Fprintf(tw, format, "---------------------------", "--------------")
	for _, i := range i {
		fmt.Fprintf(tw, format, i.Name, i.MinDisk)
	}
	tw.Flush()
}
