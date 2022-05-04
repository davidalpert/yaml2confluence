package commands

import (
	"github.com/NorthfieldIT/yaml2confluence/internal/cli"
	"github.com/NorthfieldIT/yaml2confluence/internal/services"
	"github.com/docopt/docopt-go"
)

type HooksCmd struct {
	service services.IHooksSrv
}

func (hc HooksCmd) Usage() string {
	return `
Usage:
	y2c hooks list [<instance_or_space_directory>]
	y2c hooks show <name> [<instance_or_space_directory>]

Options:
	<name> 		The name of the hook to show
`
}

func (hc HooksCmd) Handler(args docopt.Opts) {
	if args["list"].(bool) {
		dir := ToString(args["<instance_or_space_directory>"])
		if dir == "" {
			dir = "."
		}
		hc.service.List(dir)
	} else if args["show"].(bool) {
		hc.service.Show(ToString(args["<name>"]), ToString(args["<instance_or_space_directory>"]))
	}
}

func init() {
	cli.RegisterCommand("hooks", HooksCmd{services.NewHooksService()})
}
