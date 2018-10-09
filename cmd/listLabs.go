package cmd

import (
	"fmt"
	"log"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/nuagenetworks/nuxctl/nuagex"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listLabsCmd)

	listLabsCmd.Flags().StringVarP(&CredFPath, "credentials", "c", "user_creds.yml", "Path to the user credentials file.")
}

var listLabsCmd = &cobra.Command{
	Use:   "list-labs",
	Short: "Display NuageX labs",
	Long:  `Outputs to console the list of NuageX labs`,
	Run:   listLabs,
}

func listLabs(cmd *cobra.Command, args []string) {
	loginUser(cmd, args)

	fmt.Println("Retrieving NuageX Labs...")

	labs, err := nuagex.GetLabs(&user)
	if err != nil {
		log.Fatalf("%v", err)
	}

	printLabs(labs)
}

type byLabName []*nuagex.Lab

func (x byLabName) Len() int           { return len(x) }
func (x byLabName) Less(i, j int) bool { return x[i].Name < x[j].Name }
func (x byLabName) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

func printLabs(l []*nuagex.Lab) {
	const format = "%v\t%v\t%v\t%v\t%v\t%v\n"
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	sort.Sort(byLabName(l))
	fmt.Printf("\n")
	fmt.Fprintf(tw, format, "ID", "Name", "Status", "Expires", "External IP", "Password")
	fmt.Fprintf(tw, format, "------------------------", "----------------------", "-------", "----------------------", "---------------", "----------------")
	for _, l := range l {
		fmt.Fprintf(tw, format, l.ID, l.Name, l.Status, l.Expires.Format("2006-01-02 15:04 (MST)"), l.ExternalIP, l.Password)
	}
	tw.Flush()
}
