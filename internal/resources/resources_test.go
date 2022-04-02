package resources

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/NorthfieldIT/yaml2confluence/internal/utils"
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

func MockLoadYaml(paths [][]interface{}) func(string) map[interface{}]interface{} {
	yaml := make(map[string]map[interface{}]interface{})
	for _, p := range paths {
		path := p[0].(string)
		if isYamlFile(path) {
			yaml[path] = createMockYaml(p[2].(string), p[3].(string))
		}
	}
	return func(file string) map[interface{}]interface{} {
		return yaml[file]
	}
}

func createMockYaml(kind string, title string) map[interface{}]interface{} {
	if kind == "" {
		return getDefaultParentYaml(title)
	}
	return map[interface{}]interface{}{
		"kind":  kind,
		"title": title,
	}
}

func compare(t *testing.T, expected []*YamlResource, actual []*YamlResource) {
	if len(expected) != len(actual) {
		t.Fatalf("Expected exactly %d YamlResources, got %d\n%s", len(expected), len(actual), printYamlResources(expected, actual))
	}

	for i, e := range expected {
		if e.title != actual[i].title {
			t.Fatalf("For YamlResource[%d], titles do not match.\n%s", i, printYamlResources(expected, actual))
		}
		if e.path != actual[i].path {
			t.Fatalf("For YamlResource[%d], paths do not match.\n%s", i, printYamlResources(expected, actual))
		}

		//TODO compare yaml values
	}
}

func printYamlResources(e []*YamlResource, a []*YamlResource) string {
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
func TestGenerateTree(t *testing.T) {
	paths := []string{
		"/home/user/confluence/spaces/DEMO",
		"/home/user/confluence/spaces/DEMO/apps",
		"/home/user/confluence/spaces/DEMO/apps/app1.yml",
		"/home/user/confluence/spaces/DEMO/apps/nested",
		"/home/user/confluence/spaces/DEMO/apps/nested/app2.yml",
		"/home/user/confluence/spaces/DEMO/freeform.yml",
	}

	GenerateTree(paths)
}

func TestLoadYamlResources(t *testing.T) {
	paths := [][]interface{}{
		{"/home/user/confluence/spaces/DEMO", true},
		{"/home/user/confluence/spaces/DEMO/apps", true},
		{"/home/user/confluence/spaces/DEMO/apps/app1.yml", false, "application", "Test Application 1"},
		{"/home/user/confluence/spaces/DEMO/apps/index.yml", false, "wiki", "Applications"},
		{"/home/user/confluence/spaces/DEMO/apps/nested", true},
		{"/home/user/confluence/spaces/DEMO/apps/nested/app2.yml", false, "application", "Test Application 2"},
		{"/home/user/confluence/spaces/DEMO/freeform.yml", false, "wiki", "Wiki Example"},
	}

	expectedInput := [][]string{
		{"Applications", "/apps", "wiki"},
		{"Test Application 1", "/apps/app1.yml", "application"},
		{"nested", "/apps/nested", ""},
		{"Test Application 2", "/apps/nested/app2.yml", "application"},
		{"Wiki Example", "/freeform.yml", "wiki"},
	}
	expected := []*YamlResource{}
	for _, input := range expectedInput {
		expected = append(expected, &YamlResource{input[0], input[1], createMockYaml(input[2], input[0])})
	}

	walk = MockWalk(paths)
	loadYaml = MockLoadYaml(paths)

	actual := LoadYamlResources("/home/user/confluence/spaces/DEMO")

	compare(t, expected, actual)
}

func TestEnsureUniqueTitles(t *testing.T) {
	yr := []*YamlResource{
		{"test app 1", "/apps", nil},
		{"test app 2", "/apps", nil},
		{"TEST APP 2", "", nil},
	}

	err := EnsureUniqueTitles(yr)
	if err == nil {
		t.Fatal(`Expected error, got nil`)
	}

	expectedErrMsg := fmt.Sprintf(utils.DUPLICATE_TITLE, "test app 2", "/apps", "TEST APP 2", "")
	if err.Error() != expectedErrMsg {
		t.Fatalf(`Expected error of "%s", got "%s"`, expectedErrMsg, err.Error())
	}
}
