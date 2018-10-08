package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"github.com/nuagenetworks/nuxctl/nuagex"

	"gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
)

// outLabFPath is a path to the file that will receive NuageX Lab configuration dumped with `dump-lab command`
var outLabFPath string

func init() {
	rootCmd.AddCommand(dumpLabCmd)

	dumpLabCmd.Flags().StringVarP(&labID, "lab-id", "i", "", "Lab ID. Seen as the variable portion in the lab hostname.")
	dumpLabCmd.MarkFlagRequired("lab-id")

	dumpLabCmd.Flags().StringVarP(&CredFPath, "credentials", "c", "user_creds.yml", "Path to the user credentials file.")

	dumpLabCmd.Flags().StringVarP(&outLabFPath, "file", "f", "dumplab.yml", "Path to the local YAML file that will receive lab configuration.")
}

var dumpLabCmd = &cobra.Command{
	Use:   "dump-lab",
	Short: "Dump NuageX lab (environment) configuration in a file.",
	Long:  `Dump the existing NuageX lab configuration in a file.`,
	Run:   dumpLab,
}

func dumpLab(cmd *cobra.Command, args []string) {
	loginUser(cmd, args)

	fmt.Println("Retrieving NuageX Lab configuration...")

	l, _, err := nuagex.GetLab(&user, labID)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// scratch out the reason field to not appear in the YAML file
	l.Reason = ""

	// get Template from dumped lab to later exclude template entities
	// from dumped lab file
	t, _, err := nuagex.GetTemplate(&user, l.Template)
	if err != nil {
		log.Fatalf("%v", err)
	}

	ly, err := yaml.Marshal(&l)
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Printf("Writing Lab configuration to '%s' file...\n", outLabFPath)

	err = ioutil.WriteFile(outLabFPath, ly, 0644)
	if err != nil {
		log.Fatalf("%v", err)
	}
	// comment Template entities from the dumped lab file
	commentTemplateEntities(outLabFPath, t)
	fmt.Printf("Lab configuration has been successfully written to '%s' file!\n", outLabFPath)
}

func commentTemplateEntities(fp string, t *nuagex.Template) {
	tmpPath := "tmp_dumplab.yml"
	tmpFile, err := os.Create(tmpPath)
	if err != nil {
		log.Fatal(err)
	}
	defer tmpFile.Close()

	file, err := os.Open(outLabFPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	topSecFound := false // if one of the top section `services`, `servers`, or `networks:` has been found
	doComment := false
	ssnr := regexp.MustCompile(`^- name: (\S+)$`) // regexp to match a sub section name
	var sn string                                 // top section name
	var ssn string                                // sub section name
	for scanner.Scan() {
		s := scanner.Text()

		if s == "networks:" || s == "services:" || s == "servers:" {
			topSecFound = true
			sn = s[:len(s)-1]
			doComment = false // do not ever comment the top section name
			// fmt.Printf("in section %s\n", sn)
		}

		// ssm - subsection match
		ssm := ssnr.FindStringSubmatch(s)
		if topSecFound && ssm != nil && len(ssm) == 2 {
			ssn = ssm[1]
			// fmt.Printf("subsection %s\n", ssn)
			// if a pair section+subsection is found in template, we need to comment it
			if isInT(sn, ssn, t) {
				doComment = true
			} else {
				doComment = false
			}
			// fmt.Printf("is section %s and subsec %s in template?: %v\n", sn, ssn, is)
		}

		if doComment {
			s = fmt.Sprintf("# %s", s) // add comment symbol
		}

		tmpFile.WriteString(fmt.Sprintf("%s\n", s))
	}

	err = os.Rename(tmpPath, outLabFPath)
	if err != nil {
		log.Fatal(err)
	}

}

// isInT checks if a particular section+subsection pair exists in the Template `t`
func isInT(sn, ssn string, t *nuagex.Template) bool {
	var sl []string
	switch sn {
	case "services":
		{
			for _, v := range t.Services {
				sl = append(sl, v.Name)
			}
		}
	case "networks":
		{
			for _, v := range t.Networks {
				sl = append(sl, v.Name)
			}
		}
	case "servers":
		{
			for _, v := range t.Servers {
				sl = append(sl, v.Name)
			}
		}
	}
	for _, v := range sl {
		if v == ssn {
			return true
		}
	}
	return false
}
