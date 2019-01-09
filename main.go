package main

import "github.com/nuagenetworks/nuxctl/cmd"

func main() {
	var (
		Version = "0.6.0"
	)
	cmd.Execute(Version)
}
