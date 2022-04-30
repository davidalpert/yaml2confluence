package services

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/NorthfieldIT/yaml2confluence/internal/resources"
	"github.com/NorthfieldIT/yaml2confluence/internal/utils"
	"github.com/fatih/color"
)

type ITemplatesSrv interface {
	List(string)
}

type TemplatesSrv struct{}

func NewTemplatesService() ITemplatesSrv {
	return TemplatesSrv{}
}

func (TemplatesSrv) List(instanceDirectory string) {
	dirProps := utils.GetDirectoryProperties(instanceDirectory)
	tp := resources.NewTemplateProcessor(dirProps.TemplatesDir)

	assets := []resources.IAsset{}
	for _, t := range tp.GetAll() {
		assets = append(assets, t.Asset)
	}

	prettyPrintAssets(assets)
}

func prettyPrintAssets(assets []resources.IAsset) {
	writer := tabwriter.NewWriter(os.Stdout, 0, 3, 3, ' ', 0)
	bold := color.New(color.Bold)
	gray := color.New(color.FgHiBlack)

	builtin := []resources.IAsset{}
	userDefined := []resources.IAsset{}

	for _, asset := range assets {
		if asset.IsBuiltin() {
			builtin = append(builtin, asset)
		} else {
			userDefined = append(userDefined, asset)
		}
	}

	bold.Println("Built-in:")

	if len(builtin) == 0 {
		gray.Println("<none>")
	} else {
		for _, a := range builtin {
			fmt.Println(color.BlueString(a.GetName()))
		}
	}

	fmt.Println("")
	bold.Println("User-defined:")

	if len(userDefined) == 0 {
		gray.Println("<none>")
	} else {
		for _, a := range userDefined {
			fmt.Fprintln(writer, fmt.Sprintf("%s\t%s", color.BlueString(a.GetName()), a.GetPath()))
		}
	}

	writer.Flush()
}
