package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"os"
	"path/filepath"

	"github.com/NorthfieldIT/yaml2confluence/internal/hooks"
	"github.com/NorthfieldIT/yaml2confluence/internal/resources"
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
	hooks     *hooks.HookProcessor
	hasher    hash.Hash
}

func NewRenderTools(dirProps DirectoryProperties, precompileJqHooks bool) *RenderTools {
	rt := RenderTools{
		dirProps: dirProps,
		templates: map[string]string{
			"markup": "{{{markup}}}",
		},
		hooks: hooks.NewHookProcessor(dirProps.HooksDir, precompileJqHooks),
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
	hookset := rt.hooks.GetHookSet(p.Resource.Kind)

	switch {
	case target >= YAML:
		for _, yq := range hookset.Yq {
			node, err := yq.Run(p.Resource.Node)
			if err != nil {
				panic(err)
			}

			p.Resource.Node = node
		}
		p.Resource.UpdateJson()
		fallthrough
	case target >= JSON:
		for _, jq := range hookset.Jq {
			res, err := jq.Run(p.Resource.Json)
			if err != nil {
				fmt.Printf("Failed to render %s\nError in hook: %s\n\njq %s\n%s\n\n", filepath.Join(rt.dirProps.SpaceDir, p.Resource.Path), jq.Hook.Path, jq.Cmd, err.Error())
				os.Exit(1)
			}

			p.Resource.Json = res
		}
		p.Resource.UpdateKindAndTitle()
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
