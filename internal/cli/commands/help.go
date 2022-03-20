package cli

import (
	"fmt"
	"strings"

	"github.com/NorthfieldIT/yaml2confluence/internal/cli"
	"github.com/docopt/docopt-go"
)

var helpUsage = `
Usage:
	y2c help <command>
`

var helpCmd = func(args docopt.Opts) {
	fmt.Println(strings.TrimSpace(cli.GetCommandUsage(args["<command>"].(string))))
}

func init() {
	// cli.RegisterCommand("help", helpUsage, helpCmd)
}
