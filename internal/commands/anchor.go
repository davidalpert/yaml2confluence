package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/NorthfieldIT/yaml2confluence/internal/cli"
	"github.com/NorthfieldIT/yaml2confluence/internal/utils"
	"github.com/docopt/docopt-go"
)

type AnchorCmd struct {
}

func (ac AnchorCmd) Usage() string {
	return `
Usage:
	y2c anchor <instance_or_space_directory> <page_id>

Options:
	<page_id> 		The page ID to anchor the space to
`
}

func (ac AnchorCmd) Handler(args docopt.Opts) {
	dirProps := utils.GetDirectoryProperties(ToString(args["<instance_or_space_directory>"]))

	anchorFile := filepath.Join(dirProps.SpaceDir, ".anchor")

	// write .anchor file
	err := os.WriteFile(anchorFile, []byte(ToString(args["<page_id>"])), 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println("Created anchor file " + anchorFile)
}

func init() {
	cli.RegisterCommand("anchor", AnchorCmd{})
}
