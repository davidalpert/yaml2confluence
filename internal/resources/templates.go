package resources

import (
	"fmt"
	"os"
)

type TemplateProcessor struct {
	templates map[string]Template
}

type Template struct {
	Asset IAsset
	Data  string
}

func NewTemplateProcessor(templatesDir string) *TemplateProcessor {
	tp := TemplateProcessor{
		templates: loadAllTemplates(templatesDir),
	}

	return &tp
}

func (tp TemplateProcessor) Get(kind string) string {
	template, exists := tp.templates[kind]

	if !exists {
		fmt.Printf("No template exists for kind '%s'\n", kind)
		os.Exit(1)
	}

	if template.Data == "" {
		template.Data = template.Asset.ReadString()
	}

	return template.Data
}

func (tp TemplateProcessor) GetAll() []Template {
	templates := []Template{}
	for _, t := range tp.templates {
		templates = append(templates, t)
	}

	return templates
}

func loadAllTemplates(templatesDir string) map[string]Template {
	templates := map[string]Template{}

	assets := append(GetBuiltinTemplates(), LoadAssets(templatesDir, []string{".mst", ".mustache"})...)

	for _, asset := range assets {
		templates[asset.GetName()] = Template{Asset: asset}
	}

	return templates
}
