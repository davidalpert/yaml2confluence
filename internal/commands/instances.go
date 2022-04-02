package commands

import (
	"github.com/NorthfieldIT/yaml2confluence/internal/cli"
	"github.com/NorthfieldIT/yaml2confluence/internal/services"
	"github.com/docopt/docopt-go"
)

type InstanceCmd struct {
	service services.IInstancesSrv
}

func (ic InstanceCmd) Usage() string {
	return `
Usage:
	y2c instances new [<base_dir>] [--config-only]
	y2c instances list
	y2c instances decrypt <config>

Options:
	--config-only  	Skips directory creation, outputting the generated config file to the terminal
	
Sub-Commands:
	new  			A CLI wizard to configure a new Confluence instance
	list  			Lists all Confluence instances
	decrypt  		Displays a config yaml file with decrypted secrets
`
}

func (ic InstanceCmd) Handler(args docopt.Opts) {
	if args["list"].(bool) {
		ic.service.List()
	} else if args["new"].(bool) {
		ic.service.New(ToString(args["<base_dir>"]), args["--config-only"].(bool))
	}
}

//TODO move this to utils
func ToString(arg interface{}) string {
	if arg == nil {
		return ""
	}

	return arg.(string)
}

func init() {
	cli.RegisterCommand("instances", InstanceCmd{services.NewInstancesService()})
}
