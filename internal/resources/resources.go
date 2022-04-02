package resources

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/NorthfieldIT/yaml2confluence/internal/utils"
	"gopkg.in/yaml.v2"
)

type YamlResource struct {
	title string
	path  string
	yaml  map[interface{}]interface{}
}

func LoadYaml(file string) map[interface{}]interface{} {
	data, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	m := make(map[interface{}]interface{})

	err = yaml.Unmarshal(data, &m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return m
}

var loadYaml func(file string) map[interface{}]interface{} = LoadYaml

func LoadYamlResource(spaceRootDir, relFilePath string) *YamlResource {
	yaml := loadYaml(filepath.Join(spaceRootDir, relFilePath))
	if yaml["title"] == nil {
		fmt.Printf("'%s' missing required field 'title", relFilePath)
		os.Exit(1)
	}

	return &YamlResource{
		title: yaml["title"].(string),
		path:  relFilePath,
		yaml:  yaml,
	}

}

func isYamlFile(file string) bool {
	ext := filepath.Ext(file)
	return ext == ".yml" || ext == ".yaml"
}

func isIndexFile(file string) bool {
	return isYamlFile(file) && strings.Split(filepath.Base(file), ".")[0] == "index"
}

var walk func(root string, fn filepath.WalkFunc) error = filepath.Walk

func LoadYamlResources(dir string) []*YamlResource {
	resources := []*YamlResource{}
	parents := make(map[string]*YamlResource)
	dirStringLength := len(dir)

	err := walk(dir,
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
				yr := &YamlResource{
					path: relPath,
				}
				parents[relPath] = yr
				resources = append(resources, yr)
			} else if isYamlFile(path) {
				yr := LoadYamlResource(dir, relPath)
				if isIndexFile(path) {
					parent := parents[filepath.Dir(relPath)]
					parent.title = yr.title
					parent.yaml = yr.yaml
				} else {
					resources = append(resources, yr)
				}
			}

			return nil
		})
	if err != nil {
		panic(err)
	}

	populateMissingParentProperties(parents)
	return resources
}

func populateMissingParentProperties(parents map[string]*YamlResource) {
	for _, yr := range parents {
		if yr.title == "" {
			pathTokens := strings.Split(yr.path, string(os.PathSeparator))
			title := pathTokens[len(pathTokens)-1:][0]

			yr.title = title
			yr.yaml = getDefaultParentYaml(title)
		}
	}
}

func getDefaultParentYaml(title string) map[interface{}]interface{} {
	return map[interface{}]interface{}{
		"kind":   "wiki",
		"title":  title,
		"markup": "",
	}
}

func EnsureUniqueTitles(resources []*YamlResource) error {
	uniqueTitle := map[string]*YamlResource{}

	for _, cur := range resources {
		lowerTitle := strings.ToLower(cur.title)
		if r, exists := uniqueTitle[lowerTitle]; exists {
			return errors.New(fmt.Sprintf(utils.DUPLICATE_TITLE, r.title, r.path, cur.title, cur.path))
		} else {
			uniqueTitle[lowerTitle] = cur
		}
	}
	return nil
}

func GenerateTree(files []string) *utils.Node {
	if len(files) == 0 {
		return nil
	}

	root := utils.NewNode(files[0])
	files = files[1:]

	pathTokens := strings.Split(root.Value.(string), string(os.PathSeparator))
	fmt.Println(pathTokens)
	splitToken := strings.Join(pathTokens[len(pathTokens)-2:], string(os.PathSeparator))

	fmt.Println(splitToken)
	return root

}
