package commands

import (
	"fmt"
	"os"
	"testing"

	"github.com/NorthfieldIT/yaml2confluence/internal/cli"
)

type MockInstances struct {
	Calls []interface{}
}

func (mi *MockInstances) New(baseDir string, configOnly bool) {
	mi.Calls = append(mi.Calls, []interface{}{"New", baseDir, configOnly})
}

func (mi *MockInstances) List() {
	mi.Calls = append(mi.Calls, []interface{}{"List"})
}

func TestInstancesHandler(t *testing.T) {
	mockInstancesService := &MockInstances{}
	cli.RegisterCommand("instances", InstanceCmd{mockInstancesService})

	args := []string{"instances", "new", "./"}
	os.Args = os.Args[0:1]
	os.Args = append(os.Args, args...)
	cli.Parse()

	if len(mockInstancesService.Calls) != 1 {
		t.Fatalf(`Expected exactly 1 call to instances service, got %d`, len(mockInstancesService.Calls))
	}
	expectedCall := "[New ./ false]"
	actualCall := fmt.Sprint(mockInstancesService.Calls[0])

	if actualCall != expectedCall {
		t.Fatalf(`Wrong call signature, expected %s, got %s`, expectedCall, actualCall)
	}
}
