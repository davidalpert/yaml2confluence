package commands

import (
	"github.com/NorthfieldIT/yaml2confluence/internal/cli"
	"github.com/NorthfieldIT/yaml2confluence/internal/services"
	"github.com/docopt/docopt-go"
)

type TemplatesCmd struct {
	service services.ITemplatesSrv
}

func (tc TemplatesCmd) Usage() string {
	return `
Usage:
	y2c templates list [<instance_or_space_directory>]
	y2c templates show <name> [<instance_or_space_directory>] 

Options:
	<name> 		The name of the template to show
`
}

func (tc TemplatesCmd) Handler(args docopt.Opts) {
	dir := ToString(args["<instance_or_space_directory>"])
	if dir == "" {
		dir = "."
	}

	if args["list"].(bool) {
		tc.service.List(dir)
	} else if args["show"].(bool) {
		tc.service.Show(ToString(args["<name>"]), dir)
	}
}

func init() {
	cli.RegisterCommand("templates", TemplatesCmd{services.NewTemplatesService()})
}
