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
	yr := utils.LoadSingleYamlResource(file)
	page := resources.NewPage(yr.Path, yr)
	rt := utils.NewRenderTools(dirProps, true)

	target := getRenderTarget(output)
	rt.RenderTo(target, page)

	utils.PrettyPrint(target, page, os.Stdout)
}

func getRenderTarget(output string) utils.RenderTarget {
	lower := strings.ToLower(output)
	switch lower {
	case "json":
		return utils.JSON
	case "yaml":
		return utils.YAML
	default:
		return utils.MST
	}
}
