package hooks

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	. "github.com/flant/libjq-go"
	"github.com/flant/libjq-go/pkg/jq"
	"gopkg.in/yaml.v3"
)

const builtin = `target: .*
priority: 100
defaults:
  kind: markup
  labels: []
yq: 
  # Defaults 'title' to the relative path (stored in head_comment)
  - '{} as $d|$d.title = (head_comment | capture("(.*)") .[])|. *n $d'
  # Defaults 'editorVersion' to instance setting (stored in foot_comment)
  - '{} as $d|$d.editorVersion = (foot_comment | capture("(.*)") .[])|. *n $d'
  # Remove foot_comment (editorVersion)
  - '. foot_comment=""'`

type HookProcessor struct {
	shouldPrecompile bool
	kindHooks        map[string]*Hook
	patternHooks     []*Hook
}

type Hook struct {
	Name      string
	Path      string
	Target    string    `yaml:"target"`
	Priority  int       `yaml:"priority"`
	Defaults  yaml.Node `yaml:"defaults"`
	Overrides yaml.Node `yaml:"overrides"`
	Merges    yaml.Node `yaml:"merges"`
	Yq        []string  `yaml:"yq"`
	Jq        []string  `yaml:"jq"`
}

type HookSet struct {
	Jq []JqCommand
	Yq []YqHooks
}

type JqCommand struct {
	precompiled *jq.JqProgram
	Cmd         string
	Hook        *Hook
}

func (jc *JqCommand) precompile() error {
	prg, err := Jq().Program(jc.Cmd).Precompile()
	if err != nil {
		return err
	}

	jc.precompiled = prg

	return nil
}

func (jc *JqCommand) Run(json string) (string, error) {
	if jc.precompiled != nil {
		return jc.precompiled.Run(json)
	} else {
		return Jq().Program(jc.Cmd).Run(json)
	}
}

func NewHookProcessor(hooksDir string, precompile bool) *HookProcessor {
	hp := HookProcessor{
		shouldPrecompile: precompile,
		kindHooks:        map[string]*Hook{},
	}

	hooks := append(loadHooks(hooksDir), loadBuiltinHook())

	for _, hook := range hooks {
		if hook.Target == "" {
			hp.kindHooks[hook.Name] = hook
		} else {
			hp.patternHooks = append(hp.patternHooks, hook)
		}
	}

	return &hp
}

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

	sort.SliceStable(hooks, func(i, j int) bool {
		return hooks[i].Priority < hooks[j].Priority
	})

	return hooks
}

func (hp *HookProcessor) GetHookSet(kind string) HookSet {
	hookset := HookSet{}

	for _, hook := range hp.GetHooks(kind) {
		for _, jq := range hook.Jq {
			jqCommand := JqCommand{Cmd: jq, Hook: hook}
			if hp.shouldPrecompile {
				err := jqCommand.precompile()
				if err != nil {
					fmt.Printf("Failed to precompile jq statement\nHook name: %s\nFile: %s\njq: %s\nError: %s", hook.Name, hook.Path, jq, err.Error())
					os.Exit(1)
				}
			}
			hookset.Jq = append(hookset.Jq, jqCommand)
		}

		yqHooks, err := NewYqHook(hook.Defaults, hook.Overrides, hook.Merges, hook.Yq)
		if err != nil {
			panic(err)
		}

		hookset.Yq = append(hookset.Yq, yqHooks)

	}

	return hookset
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

			hook, err := loadHookFromFilePath(path)
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

func loadHookFromFilePath(path string) (*Hook, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	hook, err := loadHook(data)
	if err != nil {
		panic(err)
	}

	hook.Name = fileNameWithoutExtension(path)
	hook.Path = path

	return hook, nil
}

func loadHook(data []byte) (*Hook, error) {
	hook := Hook{}

	node := yaml.Node{}
	yaml.Unmarshal(data, &node)

	if len(node.Content) != 1 || node.Content[0].ShortTag() != "!!map" {
		return nil, errors.New("Invalid hook yaml")
	}

	ensureArray("yq", &node)
	ensureArray("jq", &node)

	err := node.Decode(&hook)
	if err != nil {
		return nil, err
	}

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

func loadBuiltinHook() *Hook {
	hook, err := loadHook([]byte(builtin))
	if err != nil {
		panic(err)
	}
	hook.Name = "builtin"

	return hook
}
