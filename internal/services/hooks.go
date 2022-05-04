package services

import (
	"os"
	"sort"

	"github.com/NorthfieldIT/yaml2confluence/internal/resources"
	"github.com/NorthfieldIT/yaml2confluence/internal/utils"
	"gopkg.in/yaml.v3"
)

type IHooksSrv interface {
	List(string)
	Show(string, string)
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

// TODO added some error handling
func (HooksSrv) Show(name, instanceDirectory string) {
	dirProps := utils.GetDirectoryProperties(instanceDirectory)
	hp := resources.NewHookProcessor(dirProps.HooksDir, false)

	node := yaml.Node{}
	yaml.Unmarshal(hp.Get(name).Asset.ReadBytes(), &node)
	resources.PrettyPrintYaml(&node, os.Stdout)
}
