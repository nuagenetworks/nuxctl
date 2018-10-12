package main

import "github.com/nuagenetworks/nuxctl/cmd"

func main() {
	var (
		Version = "0.5.1"
	)
	cmd.Execute(Version)
}
