package resources

import (
	"bytes"
	"encoding/json"
	"path/filepath"

	"github.com/aybabtme/orderedjson"
	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"gopkg.in/yaml.v3"
)

type requiredFields struct {
	Kind  string `yaml:"kind"`
	Title string `yaml:"title"`
}
type YamlResource struct {
	Kind  string
	Title string
	Path  string
	Node  *yaml.Node
	Json  string
}

func NewYamlResource(path string, node *yaml.Node) *YamlResource {
	requiredFields := &requiredFields{}
	err := node.Decode(requiredFields)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	yqlib.NewJONEncoder(0).Encode(&buf, node)

	return &YamlResource{
		Kind:  requiredFields.Kind,
		Title: requiredFields.Title,
		Path:  path,
		Node:  node,
		Json:  buf.String(),
	}
}

func (yr *YamlResource) GetParentPath() string {
	return filepath.Dir(yr.Path)
}

func (yr *YamlResource) ToObject() map[string]interface{} {
	var obj map[string]interface{}

	if err := json.Unmarshal([]byte(yr.Json), &obj); err != nil {
		panic(err)
	}

	return obj
}

func (yr *YamlResource) ToOrderedMap() orderedjson.Map {
	var object orderedjson.Map
	err := json.Unmarshal([]byte(yr.Json), &object)
	if err != nil {
		panic(err)
	}

	return object
}
