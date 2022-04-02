package services

import (
	"fmt"
	"os"

	"github.com/NorthfieldIT/yaml2confluence/internal/confluence"
	"github.com/NorthfieldIT/yaml2confluence/internal/resources"
)

type IUploadSrv interface {
	UploadSingleResource(string)
	UploadSpace(string)
}

type UploadSrv struct {
	renderSrv IRenderSrv
}

func NewUploadService() UploadSrv {
	return UploadSrv{NewRenderService()}
}

func (us UploadSrv) UploadSingleResource(file string) {
	dirProps := confluence.GetDirectoryProperties(file)
	title, markup := us.renderSrv.RenderSingleResource(file)

	confluence.CreatePage(title, markup, dirProps.SpaceKey, confluence.LoadConfig(dirProps.ConfigPath))
}

func (us UploadSrv) UploadSpace(spaceDirectory string) {
	dirProps := confluence.GetDirectoryProperties(spaceDirectory)

	yr := resources.LoadYamlResources(dirProps.SpaceDir)

	if err := resources.EnsureUniqueTitles(yr); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
