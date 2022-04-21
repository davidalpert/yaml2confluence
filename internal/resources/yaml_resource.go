package resources

import (
	"bytes"
	"path/filepath"

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
	obj := map[string]interface{}{}
	err := yr.Node.Decode(&obj)
	if err != nil {
		panic(err)
	}

	return obj
}
