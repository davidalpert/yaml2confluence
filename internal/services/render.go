package services

import (
	"github.com/NorthfieldIT/yaml2confluence/internal/resources"
	"github.com/NorthfieldIT/yaml2confluence/internal/utils"
)

type IRenderSrv interface {
	RenderSingleResource(string) string
}

type RenderSrv struct{}

func NewRenderService() RenderSrv {
	return RenderSrv{}
}

func (RenderSrv) RenderSingleResource(file string) string {
	dirProps := utils.GetDirectoryProperties(file)
	yr := utils.LoadSingleYamlResource(file)

	utils.RenderPage(resources.NewPage(yr.Path, yr), utils.LoadTemplate(yr.Kind, dirProps.TemplatesDir), utils.LoadHook(yr.Kind, dirProps.HooksDir))

	return yr.Json
}
