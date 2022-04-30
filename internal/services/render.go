package services

import (
	"os"
	"strings"

	"github.com/NorthfieldIT/yaml2confluence/internal/resources"
	"github.com/NorthfieldIT/yaml2confluence/internal/utils"
)

type IRenderSrv interface {
	RenderSingleResource(string, string)
}

type RenderSrv struct{}

func NewRenderService() RenderSrv {
	return RenderSrv{}
}

func (RenderSrv) RenderSingleResource(file string, output string) {
	dirProps := utils.GetDirectoryProperties(file)
	yr := resources.LoadSingleYamlResource(file)
	page := resources.NewPage(yr.Path, yr)
	rt := resources.NewRenderTools(dirProps, true)

	target := getRenderTarget(output)
	rt.RenderTo(target, page)

	resources.PrettyPrint(target, page, os.Stdout)
}

func getRenderTarget(output string) resources.RenderTarget {
	lower := strings.ToLower(output)
	switch lower {
	case "json":
		return resources.JSON
	case "yaml":
		return resources.YAML
	default:
		return resources.MST
	}
}
