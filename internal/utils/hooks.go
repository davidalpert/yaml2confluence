package utils

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type HookProcessor struct {
	kindHooks    map[string]*Hook
	patternHooks []*Hook
}

type Hook struct {
	name   string
	path   string
	Target string `yaml:"target"`
	Jq     []string
}

func NewHookProcessor(hooksDir string, precompile bool) *HookProcessor {
	hp := HookProcessor{
		kindHooks: map[string]*Hook{},
	}

	hooks := loadHooks(hooksDir)

	for _, hook := range hooks {
		if hook.Target == "" {
			hp.kindHooks[hook.name] = hook
		} else {
			hp.patternHooks = append(hp.patternHooks, hook)
		}
	}

	return &hp
}

// really should be returning a slice of a specific type of hook i.e GetJqHooks()
func (hp *HookProcessor) GetHooks(kind string) []*Hook {
	hooks := []*Hook{}

	if kindHook, exists := hp.kindHooks[kind]; exists {
		hooks = append(hooks, kindHook)
	}

	for _, ph := range hp.patternHooks {
		if matched, _ := regexp.MatchString(ph.Target, kind); matched {
			hooks = append(hooks, ph)
		}
	}
	return hooks
}

func loadHooks(hooksDir string) []*Hook {
	hooks := []*Hook{}
	err := filepath.Walk(hooksDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				panic(err)
			}

			if info.IsDir() {
				return nil
			}

			hook, err := loadHook(path)
			if err != nil {
				panic(err)
			}
			hooks = append(hooks, hook)
			return nil
		})
	if err != nil {
		panic(err)
	}

	return hooks
}

func loadHook(path string) (*Hook, error) {
	hook := Hook{}

	data, _ := os.ReadFile(path)
	node := yaml.Node{}
	yaml.Unmarshal(data, &node)

	if len(node.Content) != 1 || node.Content[0].ShortTag() != "!!map" {
		return nil, errors.New("Invalid hook yaml")
	}

	ensureArray("jq", &node)

	err := node.Decode(&hook)
	if err != nil {
		return nil, err
	}

	hook.name = fileNameWithoutExtension(path)
	hook.path = path

	return &hook, nil
}

/*
allows hooks to be defined as single value or an array

jq: .user as $s | .user |= "mike"

BECOMES

jq:
	- .user as $s | .user |= "mike"
*/
func ensureArray(rootKey string, node *yaml.Node) {
	content := node.Content[0].Content
	for i := range content {
		if content[i].Value == rootKey && content[i+1].ShortTag() == "!!str" {
			seq := yaml.Node{
				Kind:    yaml.SequenceNode,
				Content: append([]*yaml.Node{}, content[i+1]),
			}
			content[i+1] = &seq
			break
		}
	}
}

func fileNameWithoutExtension(path string) string {
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}
