package commands

import (
	"github.com/NorthfieldIT/yaml2confluence/internal/cli"
	"github.com/NorthfieldIT/yaml2confluence/internal/services"
	"github.com/docopt/docopt-go"
)

type UploadCmd struct {
	service services.IUploadSrv
}

func (ic UploadCmd) Usage() string {
	return `
Usage:
	y2c upload <space_directory>
	y2c upload -f <file> | --file <file>

Options:
	-f <file>, --file <file>     	The YAML resource to upload
`
}

func (ic UploadCmd) Handler(args docopt.Opts) {
	if spaceDir := ToString(args["<space_directory>"]); spaceDir != "" {
		ic.service.UploadSpace(spaceDir)
	} else if file := ToString(args["--file"]); file != "" {
		ic.service.UploadSingleResource(args["--file"].(string))
	}
}

func init() {
	cli.RegisterCommand("upload", UploadCmd{services.NewUploadService()})
}
