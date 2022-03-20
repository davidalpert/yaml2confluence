package cli

import (
	"fmt"

	"github.com/NorthfieldIT/yaml2confluence/internal/cli"
	"github.com/NorthfieldIT/yaml2confluence/internal/utils"
	"github.com/docopt/docopt-go"
)

var renderUsage = `
Usage:
	y2c render -f <file> | --file <file>

Options:
	-f <file>, --file <file>     	The YAML resource to render
`

var renderCmd = func(args docopt.Opts) {
	yaml := utils.LoadYaml(args["--file"].(string))
	fmt.Println(utils.RenderTemplate(yaml))
}

func init() {
	cli.RegisterCommand("render", renderUsage, renderCmd)
}
