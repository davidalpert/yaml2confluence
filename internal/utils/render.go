package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"os"
	"path/filepath"

	"github.com/NorthfieldIT/yaml2confluence/internal/resources"
	. "github.com/flant/libjq-go"
	"github.com/hoisie/mustache"
)

type RenderTarget uint32

const (
	YAML = 1 << iota
	JSON
	MST
)

type RenderTools struct {
	dirProps  DirectoryProperties
	templates map[string]string
	hooks     *HookProcessor
	hasher    hash.Hash
}

func NewRenderTools(dirProps DirectoryProperties, precompileJqHooks bool) *RenderTools {
	rt := RenderTools{
		dirProps:  dirProps,
		templates: map[string]string{},
		hooks:     NewHookProcessor(dirProps.HooksDir, precompileJqHooks),
	}

	return &rt
}

func (rt *RenderTools) GetTemplate(kind string) string {
	template, exists := rt.templates[kind]
	if !exists {
		template = loadTemplate(kind, rt.dirProps.TemplatesDir)
		rt.templates[kind] = template
	}

	return template
}

func (rt *RenderTools) RenderTo(target RenderTarget, p *resources.Page) {
	hooks := rt.hooks.GetHooks(p.Resource.Kind)

	switch {
	case target >= JSON:
		for _, hook := range hooks {
			for _, jq := range hook.Jq {
				res, err := Jq().Program(jq).Run(p.Resource.Json)
				if err != nil {
					fmt.Printf("Failed to render %s\nError in hook: %s\n\njq %s\n%s\n\n", filepath.Join(rt.dirProps.SpaceDir, p.Resource.Path), hook.path, jq, err.Error())
					os.Exit(1)
				}

				p.Resource.Json = res
			}

		}
		fallthrough
	case target == MST:
		renderContent(p, rt.GetTemplate(p.Resource.Kind))
	}
}

func (rt *RenderTools) RenderAll(pt *resources.PageTree) {
	for _, page := range pt.GetPages() {
		rt.RenderTo(MST, page)
	}
}

func renderContent(p *resources.Page, template string) {
	p.Content.Markup = mustache.Render(template, p.Resource.ToObject())
	hasher := sha256.New()
	hasher.Write([]byte(p.Content.Markup))
	p.Content.Sha256 = hex.EncodeToString(hasher.Sum(nil))
}

func loadTemplate(kind, templatesDir string) string {
	data, err := os.ReadFile(filepath.Join(templatesDir, kind+".mst"))
	if err != nil {
		panic(err)
	}

	return string(data)
}
