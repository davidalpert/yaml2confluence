package utils

import (
	"os"

	"github.com/hoisie/mustache"
)

//TODO needs to be discoverable by convention
var TEMPLATE_DIR = "instances/yaml2confluence/templates"

func RenderTemplate(data map[interface{}]interface{}) string {
	kind := data["kind"].(string)
	template := loadTemplate(kind)

	return mustache.Render(template, data)
}

func loadTemplate(kind string) string {
	data, err := os.ReadFile(TEMPLATE_DIR + "/" + kind + ".mst")
	if err != nil {
		panic(err)
	}

	return string(data)
}
