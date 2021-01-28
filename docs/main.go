package main

import (
	"log"

	"github.com/infamousjoeg/cybr-cli/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	err := doc.GenMarkdownTree(cmd.GetCMD(), "./")
	if err != nil {
		log.Fatalf("Failed to generate markdown documentation")
	}
}
