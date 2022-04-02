package commands

import (
	"fmt"

	"github.com/NorthfieldIT/yaml2confluence/internal/cli"
	"github.com/NorthfieldIT/yaml2confluence/internal/services"
	"github.com/docopt/docopt-go"
)

type RenderCmd struct {
	service services.IRenderSrv
}

func (RenderCmd) Usage() string {
	return `
Usage:
	y2c render -f <file> | --file <file>

Options:
	-f <file>, --file <file>     	The YAML resource to render
`
}

func (rc RenderCmd) Handler(args docopt.Opts) {
	_, markup := rc.service.RenderSingleResource(args["--file"].(string))
	fmt.Println(markup)
}

func init() {
	cli.RegisterCommand("render", RenderCmd{services.NewRenderService()})
}
