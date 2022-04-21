package resources

import (
	"fmt"
	"strings"
	"testing"
)

func validateLevels(expected, actual [][]string, t *testing.T) {
	if len(expected) != len(actual) {
		t.Fatalf("Number of levels do not match.\n\texpected:\n\t\t%d\n\tactual:\n\t\t%d", len(expected), len(actual))
	}
	for i := range expected {
		if len(expected[i]) != len(actual[i]) {
			t.Fatalf("Level %d length wrong.\n\texpected:\n\t\t%s\n\tactual:\n\t\t%s", i, fmt.Sprint(expected), fmt.Sprint(actual))
		}

		keys := map[string]bool{}
		for _, key := range expected[i] {
			keys[key] = true
		}

		for _, key := range actual[i] {
			if !keys[key] {
				t.Fatalf("Did not expect '%s' in Level %d.\n\texpected:\n\t\t%s\n\tactual:\n\t\t%s", key, i, fmt.Sprint(expected), fmt.Sprint(actual))
			}
		}
	}
}

func TestNewPageTree(t *testing.T) {
	yr := []*YamlResource{
		{"wiki", "apps root", "/apps", nil, ""},
		{"wiki", "test app 1", "/apps/app1.yml", nil, ""},
		{"wiki", "test app 2", "/apps/app2.yml", nil, ""},
		{"wiki", "test app 3", "/apps/app3.yml", nil, ""},
		{"wiki", "nested apps root", "/apps/nested", nil, ""},
		{"wiki", "test app 4", "/apps/nested/app2.yml", nil, ""},
		{"wiki", "test app 5", "/apps/nested/app3.yml", nil, ""},
	}

	pt := NewPageTree(yr)

	expectedLevels := [][]string{
		{"/apps"},
		{"/apps/app1.yml", "/apps/app2.yml", "/apps/app3.yml", "/apps/nested"},
		{"/apps/nested/app2.yml", "/apps/nested/app3.yml"},
	}
	actualLevels := pt.GetLevels()

	validateLevels(expectedLevels, actualLevels, t)

	page := pt.GetPageFromTitlePath([]string{"apps root", "nested apps root", "test app 4"})
	fmt.Println(page)

}
func generatePageUpdate(op ChangeType, id string) PageUpdate {
	return PageUpdate{
		Operation: op,
		Page: &Page{
			Remote: &RemoteResource{
				Id: id,
			},
		},
	}
}

func pageUpdateToString(pu [][]PageUpdate) string {
	ids := []string{}
	for _, updates := range pu {
		idStr := ""
		for _, update := range updates {
			idStr += update.Page.Remote.Id
		}
		ids = append(ids, idStr)
	}

	return strings.Join(ids, "|")
}
func TestMergePageUpdates(t *testing.T) {
	pu1 := [][]PageUpdate{
		{
			generatePageUpdate(DELETE, "1"),
			generatePageUpdate(DELETE, "2"),
		},
		{
			generatePageUpdate(DELETE, "3"),
			generatePageUpdate(DELETE, "4"),
			generatePageUpdate(DELETE, "5"),
		},
	}
	pu2 := [][]PageUpdate{
		{
			generatePageUpdate(UPDATE, "6"),
			generatePageUpdate(UPDATE, "7"),
			generatePageUpdate(DELETE, "8"),
			generatePageUpdate(DELETE, "9"),
		},
	}

	m := mergePageUpdates(pu1, pu2)
	mergeString := pageUpdateToString(m)

	if mergeString != "126789|345" {
		t.Fatalf(`Merge incorrect. Expected: "126789|345" Actual: "%s"`, mergeString)
	}

	m = mergePageUpdates(pu2, pu1)
	mergeString = pageUpdateToString(m)

	if mergeString != "678912|345" {
		t.Fatalf(`Merge incorrect. Expected: "126789|345" Actual: "%s"`, mergeString)
	}

	m = mergePageUpdates([][]PageUpdate{}, pu2)
	mergeString = pageUpdateToString(m)

	if mergeString != "6789" {
		t.Fatalf(`Merge incorrect. Expected: "126789|345" Actual: "%s"`, mergeString)
	}

	m = mergePageUpdates(pu2, [][]PageUpdate{})
	mergeString = pageUpdateToString(m)

	if mergeString != "6789" {
		t.Fatalf(`Merge incorrect. Expected: "126789|345" Actual: "%s"`, mergeString)
	}

}
