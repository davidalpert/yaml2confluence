package resources

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"

	. "github.com/flant/libjq-go"
	"github.com/flant/libjq-go/pkg/jq"
	"gopkg.in/yaml.v3"
)

type HookProcessor struct {
	shouldPrecompile bool
	hooks            map[string]*Hook
	kindHooks        map[string]*Hook
	patternHooks     []*Hook
}

type Hook struct {
	Asset  IAsset
	Config *HookConfig
}

type HookConfig struct {
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
		hooks:            map[string]*Hook{},
		kindHooks:        map[string]*Hook{},
	}

	hooks := append(loadHooks(hooksDir))

	for _, hook := range hooks {
		hp.hooks[hook.Asset.GetName()] = hook
		if hook.Config.Target == "" {
			hp.kindHooks[hook.Asset.GetName()] = hook
		} else {
			hp.patternHooks = append(hp.patternHooks, hook)
		}
	}

	return &hp
}

func (hp *HookProcessor) Get(hookName string) *Hook {
	return hp.hooks[hookName]
}

func (hp *HookProcessor) GetHooks(kind string) []*Hook {
	hooks := []*Hook{}

	if kindHook, exists := hp.kindHooks[kind]; exists {
		hooks = append(hooks, kindHook)
	}

	for _, ph := range hp.patternHooks {
		if matched, _ := regexp.MatchString(ph.Config.Target, kind); matched {
			hooks = append(hooks, ph)
		}
	}

	sort.SliceStable(hooks, func(i, j int) bool {
		return hooks[i].Config.Priority < hooks[j].Config.Priority
	})

	return hooks
}
func (hp *HookProcessor) GetAll() []*Hook {
	hooks := append([]*Hook{}, hp.patternHooks...)

	for _, h := range hp.kindHooks {
		hooks = append(hooks, h)
	}

	return hooks
}

func (hp *HookProcessor) GetHookSet(kind string) HookSet {
	hookset := HookSet{}

	for _, hook := range hp.GetHooks(kind) {
		for _, jq := range hook.Config.Jq {
			jqCommand := JqCommand{Cmd: jq, Hook: hook}
			if hp.shouldPrecompile {
				err := jqCommand.precompile()
				if err != nil {
					fmt.Printf("Failed to precompile jq statement\nHook name: %s\nFile: %s\njq: %s\nError: %s", hook.Asset.GetName(), hook.Asset.GetPath(), jq, err.Error())
					os.Exit(1)
				}
			}
			hookset.Jq = append(hookset.Jq, jqCommand)
		}

		yqHooks, err := NewYqHook(hook.Config.Defaults, hook.Config.Overrides, hook.Config.Merges, hook.Config.Yq)
		if err != nil {
			panic(err)
		}

		hookset.Yq = append(hookset.Yq, yqHooks)

	}

	return hookset
}

func loadHooks(hooksDir string) []*Hook {
	hooks := []*Hook{}

	assets := append(GetBuiltinHooks(), LoadAssets(hooksDir, []string{".yml", ".yaml"}, true)...)

	for _, asset := range assets {
		config, err := loadHookConfig(asset.ReadBytes())
		if err != nil {
			panic(err)
		}

		hook := Hook{
			Asset:  asset,
			Config: config,
		}

		hooks = append(hooks, &hook)
	}

	return hooks
}

func loadHookConfig(data []byte) (*HookConfig, error) {
	hookConfig := HookConfig{}

	node := yaml.Node{}
	yaml.Unmarshal(data, &node)

	if len(node.Content) != 1 || node.Content[0].ShortTag() != "!!map" {
		return nil, errors.New("Invalid hook yaml")
	}

	ensureArray("yq", &node)
	ensureArray("jq", &node)

	err := node.Decode(&hookConfig)
	if err != nil {
		return nil, err
	}

	return &hookConfig, nil
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
