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
	fmt.Println("NEW!!")
	mi.Calls = append(mi.Calls, []interface{}{"New", baseDir, configOnly})
	fmt.Println(mi.Calls)
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

	t.Fatalf(fmt.Sprint(mockInstancesService.Calls))
}
