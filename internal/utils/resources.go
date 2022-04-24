package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/NorthfieldIT/yaml2confluence/internal/resources"
	"gopkg.in/yaml.v3"
)

type YamlResourceLoader struct {
	Walk     func(root string, fn filepath.WalkFunc) error
	LoadYaml func(file string) []byte
}

func DefaultYamlResourceLoader() YamlResourceLoader {
	return YamlResourceLoader{filepath.Walk, DefaultLoadYaml}
}
func LoadYamlResources(dir string) []*resources.YamlResource {
	return DefaultYamlResourceLoader().loadYamlResources(dir)
}

func DefaultLoadYaml(file string) []byte {
	data, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	return data
}

func unmarshal(data []byte) *yaml.Node {
	node := yaml.Node{}
	err := yaml.Unmarshal(data, &node)
	if err != nil {
		panic(err)
	}

	return &node
}

func (yrl YamlResourceLoader) loadYamlResources(dir string) []*resources.YamlResource {
	yrs := []*resources.YamlResource{}
	parents := map[string]*resources.YamlResource{}
	dirStringLength := len(dir)

	err := yrl.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				panic(err)
			}

			// skip space dir
			if path == dir {
				return nil
			}

			relPath := path[dirStringLength:]

			if info.IsDir() {
				yr := getDefaultDirYamlResource(relPath)

				// save a pointer to the directory YamlResource for later in case an index.yml is found
				parents[relPath] = yr
				yrs = append(yrs, yr)
			} else if IsYamlFile(path) {
				yr := yrl.LoadYamlResource(dir, relPath)
				if isIndexFile(path) {
					parent := parents[filepath.Dir(relPath)]
					parent.Kind = yr.Kind
					parent.Title = yr.Title
					parent.Json = yr.Json
					parent.Node = yr.Node
				} else {
					yrs = append(yrs, yr)
				}
			}

			return nil
		})
	if err != nil {
		panic(err)
	}

	return yrs
}

func LoadSingleYamlResource(file string) *resources.YamlResource {
	fileAbs := ResolveAbsolutePathFile(file)
	yrl := YamlResourceLoader{func(root string, fn filepath.WalkFunc) error {
		file, _ := os.Stat(fileAbs)
		fn(fileAbs, file, nil)
		return nil
	}, DefaultLoadYaml}

	return yrl.loadYamlResources(GetDirectoryProperties(file).SpaceDir)[0]
}

func (yrl YamlResourceLoader) LoadYamlResource(spaceRootDir, relFilePath string) *resources.YamlResource {
	return resources.NewYamlResource(relFilePath, unmarshal(yrl.LoadYaml(filepath.Join(spaceRootDir, relFilePath))))
}

func IsYamlFile(file string) bool {
	ext := filepath.Ext(file)
	return ext == ".yml" || ext == ".yaml"
}

func isIndexFile(file string) bool {
	name := strings.Split(filepath.Base(file), ".")[0]
	return IsYamlFile(file) && (name == "index" || name == "_index")
}

func getDefaultDirYamlResource(relPath string) *resources.YamlResource {
	pathTokens := strings.Split(relPath, string(os.PathSeparator))
	title := pathTokens[len(pathTokens)-1:][0]

	return resources.NewYamlResource(relPath, unmarshal([]byte(fmt.Sprintf("kind: wiki\ntitle: %s\nmarkup: \"\"", title))))
}

func EnsureRequiredFieldsAndUniqueTitles(yrs []*resources.YamlResource) error {
	uniqueTitle := map[string]*resources.YamlResource{}

	for _, cur := range yrs {
		if cur.Title == "" {
			return errors.New(fmt.Sprintf(MISSING_FIELD, cur.Path, "title"))
		}
		if cur.Kind == "" {
			return errors.New(fmt.Sprintf(MISSING_FIELD, cur.Path, "kind"))
		}

		lowerTitle := strings.ToLower(cur.Title)
		if r, exists := uniqueTitle[lowerTitle]; exists {
			return errors.New(fmt.Sprintf(DUPLICATE_TITLE, r.Title, r.Path, cur.Title, cur.Path))
		} else {
			uniqueTitle[lowerTitle] = cur
		}
	}
	return nil
}
