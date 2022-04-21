package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"

	"github.com/NorthfieldIT/yaml2confluence/internal/resources"
	. "github.com/flant/libjq-go"
	"github.com/hoisie/mustache"
	"gopkg.in/yaml.v3"
)

type Hook struct {
	Jq string `yaml:"jq"`
}

func RenderAll(pt *resources.PageTree, templatesDir string, hooksDir string) {
	templates := map[string]string{}
	hooks := map[string]*Hook{}
	for _, page := range pt.GetPages() {
		template, exists := templates[page.Resource.Kind]
		if !exists {
			template = LoadTemplate(page.Resource.Kind, templatesDir)
			templates[page.Resource.Kind] = template
		}

		hook, exists := hooks[page.Resource.Kind]
		if !exists {
			hook = LoadHook(page.Resource.Kind, hooksDir)
			hooks[page.Resource.Kind] = hook
		}

		RenderPage(page, template, hook)
	}

}

func RenderPage(p *resources.Page, template string, hook *Hook) {
	if hook != nil && hook.Jq != "" {
		res, err := Jq().Program(hook.Jq).Run(p.Resource.Json)
		if err == nil {
			p.Resource.Json = res
		}
	}

	p.Content.Markup = mustache.Render(template, p.Resource.ToObject())
	hasher := sha256.New()
	hasher.Write([]byte(p.Content.Markup))
	p.Content.Sha256 = hex.EncodeToString(hasher.Sum(nil))
}

func LoadTemplate(kind, templatesDir string) string {
	data, err := os.ReadFile(filepath.Join(templatesDir, kind+".mst"))
	if err != nil {
		panic(err)
	}

	return string(data)
}

func LoadHook(kind, hooksDir string) *Hook {
	hook := Hook{}

	data, err := os.ReadFile(filepath.Join(hooksDir, kind+".yml"))
	if err != nil {
		return nil
	}

	err = yaml.Unmarshal([]byte(data), &hook)
	if err != nil {
		return nil
	}

	return &hook
}
