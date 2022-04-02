package services

import (
	"github.com/NorthfieldIT/yaml2confluence/internal/resources"
	"github.com/NorthfieldIT/yaml2confluence/internal/utils"
)

type IRenderSrv interface {
	RenderSingleResource(string) (string, string)
}

type RenderSrv struct{}

func NewRenderService() RenderSrv {
	return RenderSrv{}
}

func (RenderSrv) RenderSingleResource(file string) (string, string) {
	yaml := resources.LoadYaml(file)

	return yaml["title"].(string), utils.RenderTemplate(yaml)
}
