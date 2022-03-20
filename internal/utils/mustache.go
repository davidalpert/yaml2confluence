package utils

import (
	"os"

	"github.com/hoisie/mustache"
)

var TEMPLATE_DIR = "./resources/templates"

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
