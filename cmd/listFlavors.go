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
	rootCmd.AddCommand(listFlavorsCmd)

	listFlavorsCmd.Flags().StringVarP(&CredFPath, "credentials", "c", "user_creds.yml", "Path to the user credentials file.")
}

var listFlavorsCmd = &cobra.Command{
	Use:   "list-flavors",
	Short: "Display NuageX image flavors.",
	Long:  `Outputs to console the list of available NuageX flavors.`,
	Run:   listFlavors,
}

func listFlavors(cmd *cobra.Command, args []string) {
	loginUser(cmd, args)

	fmt.Println("Retrieving NuageX Lab flavors...")

	flavors, err := nuagex.GetFlavors(&user)
	if err != nil {
		log.Fatalf("%v", err)
	}

	printFlavors(flavors)
}

type byFlavorName []*nuagex.Flavor

func (x byFlavorName) Len() int           { return len(x) }
func (x byFlavorName) Less(i, j int) bool { return x[i].Name < x[j].Name }
func (x byFlavorName) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

func printFlavors(f []*nuagex.Flavor) {
	const format = "%v\t%v\t%v\t%v\n"
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	sort.Sort(byFlavorName(f))
	fmt.Printf("\n")
	fmt.Fprintf(tw, format, "Name", "CPU", "Memory (GB)", "Disk (GB)")
	fmt.Fprintf(tw, format, "------------", "---", "-----------", "---------")
	for _, f := range f {
		fmt.Fprintf(tw, format, f.Name, f.CPU, f.Memory, f.Disk)
	}
	tw.Flush()
}
