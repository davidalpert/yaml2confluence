package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/NorthfieldIT/yaml2confluence/internal/resources"
	"gopkg.in/yaml.v3"
)

type MockFileInfo struct {
	isDir bool
}

func (MockFileInfo) Name() string       { return "" }
func (MockFileInfo) Size() int64        { return 0 }
func (MockFileInfo) Mode() os.FileMode  { return 0 }
func (MockFileInfo) ModTime() time.Time { return time.Now() }
func (mfi MockFileInfo) IsDir() bool    { return mfi.isDir }
func (MockFileInfo) Sys() any           { return nil }

func MockWalk(paths [][]interface{}) func(string, filepath.WalkFunc) error {
	return func(root string, fn filepath.WalkFunc) error {
		for _, p := range paths {
			fn(p[0].(string), MockFileInfo{p[1].(bool)}, nil)
		}

		return nil
	}
}

func MockLoadYamlFileAsJson(paths [][]interface{}) func(string) []byte {
	jsonData := make(map[string][]byte)
	for _, p := range paths {
		path := p[0].(string)
		if IsYamlFile(path) {
			jsonData[path] = []byte(fmt.Sprintf("kind: %s\ntitle: %s", p[2].(string), p[3].(string)))
		}
	}
	return func(file string) []byte {
		return jsonData[file]
	}
}

func createYamlNode(kind string, title string) *yaml.Node {
	return unmarshal([]byte(fmt.Sprintf("kind: %s\ntitle: %s", kind, title)))
}

func compare(t *testing.T, expected []*resources.YamlResource, actual []*resources.YamlResource) {
	if len(expected) != len(actual) {
		t.Fatalf("Expected exactly %d YamlResources, got %d\n%s", len(expected), len(actual), printYamlResources(expected, actual))
	}

	for i, e := range expected {
		if e.Kind != actual[i].Kind {
			t.Fatalf("For YamlResource[%d], kinds do not match.\n%s", i, printYamlResources(expected, actual))
		}
		if e.Title != actual[i].Title {
			t.Fatalf("For YamlResource[%d], titles do not match.\n%s", i, printYamlResources(expected, actual))
		}
		if e.Path != actual[i].Path {
			t.Fatalf("For YamlResource[%d], paths do not match.\n%s", i, printYamlResources(expected, actual))
		}

		//TODO compare yaml values
	}
}

func printYamlResources(e []*resources.YamlResource, a []*resources.YamlResource) string {
	output := "\texpected:\n"
	for _, yr := range e {
		output = output + "\t\t" + fmt.Sprint(yr) + "\n"
	}
	output = output + "\tactual:\n"
	for _, yr := range a {
		output = output + "\t\t" + fmt.Sprint(yr) + "\n"
	}
	return output
}

func TestLoadYamlResources(t *testing.T) {
	paths := [][]interface{}{
		{"/home/user/confluence/spaces/DEMO", true},
		{"/home/user/confluence/spaces/DEMO/apps", true},
		{"/home/user/confluence/spaces/DEMO/apps/app1.yml", false, "application", "Test Application 1"},
		{"/home/user/confluence/spaces/DEMO/apps/index.yml", false, "index", "Applications"},
		{"/home/user/confluence/spaces/DEMO/apps/nested", true},
		{"/home/user/confluence/spaces/DEMO/apps/nested/app2.yml", false, "application", "Test Application 2"},
		{"/home/user/confluence/spaces/DEMO/freeform.yml", false, "wiki", "Wiki Example"},
	}

	expected := []*resources.YamlResource{
		resources.NewYamlResource("/apps", createYamlNode("index", "Applications")),
		resources.NewYamlResource("/apps/app1.yml", createYamlNode("application", "Test Application 1")),
		resources.NewYamlResource("/apps/nested", createYamlNode("wiki", "nested")),
		resources.NewYamlResource("/apps/nested/app2.yml", createYamlNode("application", "Test Application 2")),
		resources.NewYamlResource("/freeform.yml", createYamlNode("wiki", "Wiki Example")),
	}

	yrl := YamlResourceLoader{
		MockWalk(paths),
		MockLoadYamlFileAsJson(paths),
	}

	actual := yrl.loadYamlResources("/home/user/confluence/spaces/DEMO")

	compare(t, expected, actual)
}

func TestEnsureRequiredFieldsAndUniqueTitles(t *testing.T) {
	// test duplicate title
	yr := []*resources.YamlResource{
		resources.NewYamlResource("/app1.yml", createYamlNode("wiki", "test app 1")),
		resources.NewYamlResource("/app2.yml", createYamlNode("wiki", "test app 2")),
		resources.NewYamlResource("/app3.yml", createYamlNode("wiki", "TEST APP 2")),
	}

	err := EnsureRequiredFieldsAndUniqueTitles(yr)
	if err == nil {
		t.Fatal(`Expected error, got nil`)
	}

	expectedErrMsg := fmt.Sprintf(DUPLICATE_TITLE, "test app 2", "/app2.yml", "TEST APP 2", "/app3.yml")
	if err.Error() != expectedErrMsg {
		t.Fatalf(`Expected error of "%s", got "%s"`, expectedErrMsg, err.Error())
	}

	resources.NewYamlResource("/app.yml", createYamlNode("", "app"))

	// test kind is empty string
	err = EnsureRequiredFieldsAndUniqueTitles([]*resources.YamlResource{resources.NewYamlResource("/app.yml", createYamlNode("", "app"))})
	if err == nil {
		t.Fatal(`Expected error, got nil`)
	}

	expectedErrMsg = fmt.Sprintf(MISSING_FIELD, "/app.yml", "kind")
	if err.Error() != expectedErrMsg {
		t.Fatalf(`Expected error of "%s", got "%s"`, expectedErrMsg, err.Error())
	}

	resources.NewYamlResource("/app.yml", createYamlNode("wiki", ""))
	// test title is empty string
	err = EnsureRequiredFieldsAndUniqueTitles([]*resources.YamlResource{resources.NewYamlResource("/app.yml", createYamlNode("wiki", ""))})
	if err == nil {
		t.Fatal(`Expected error, got nil`)
	}

	expectedErrMsg = fmt.Sprintf(MISSING_FIELD, "/app.yml", "title")
	if err.Error() != expectedErrMsg {
		t.Fatalf(`Expected error of "%s", got "%s"`, expectedErrMsg, err.Error())
	}
}
