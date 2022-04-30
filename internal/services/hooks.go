package services

import (
	"sort"

	"github.com/NorthfieldIT/yaml2confluence/internal/resources"
	"github.com/NorthfieldIT/yaml2confluence/internal/utils"
)

type IHooksSrv interface {
	List(string)
}

type HooksSrv struct{}

func NewHooksService() HooksSrv {
	return HooksSrv{}
}

func (HooksSrv) List(instanceDirectory string) {
	dirProps := utils.GetDirectoryProperties(instanceDirectory)
	hp := resources.NewHookProcessor(dirProps.HooksDir, false)

	assets := []resources.IAsset{}
	for _, h := range hp.GetAll() {
		assets = append(assets, h.Asset)
	}
	sort.SliceStable(assets, func(i, j int) bool {
		return assets[i].GetName() < assets[j].GetName()
	})

	prettyPrintAssets(assets)
}
