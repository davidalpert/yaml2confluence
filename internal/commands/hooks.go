package commands

import (
	"github.com/NorthfieldIT/yaml2confluence/internal/cli"
	"github.com/NorthfieldIT/yaml2confluence/internal/services"
	"github.com/docopt/docopt-go"
)

type HooksCmd struct {
	service services.IHooksSrv
}

func (ic HooksCmd) Usage() string {
	return `
Usage:
	y2c hooks list [<instance_or_space_directory>]
	y2c hooks show [<instance_or_space_directory>] <name>

Options:
	<name> 		The name of the hook to show
`
}

func (ic HooksCmd) Handler(args docopt.Opts) {
	if args["list"].(bool) {
		dir := ToString(args["<instance_or_space_directory>"])
		if dir == "" {
			dir = "."
		}
		ic.service.List(dir)
	}
	// if spaceDir := ToString(args["<space_directory>"]); spaceDir != "" {
	// 	ic.service.UploadSpace(spaceDir)
	// } else if file := ToString(args["--file"]); file != "" {
	// 	ic.service.UploadSingleResource(args["--file"].(string))
	// }
}

func init() {
	cli.RegisterCommand("hooks", HooksCmd{services.NewHooksService()})
}
