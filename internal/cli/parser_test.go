package cli

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/NorthfieldIT/yaml2confluence/internal/constants"
	"github.com/docopt/docopt-go"
)

type MockInstancesCmd struct {
	UsageStr string
	Calls    int
}

func (mic *MockInstancesCmd) Usage() string       { return mic.UsageStr }
func (mic *MockInstancesCmd) Handler(docopt.Opts) { mic.Calls++ }

func ExecParseWithMocks(command string, args string) (string, int) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	var exitCode int
	exitFunction = func(code int) {
		exitCode = code
	}

	os.Args = os.Args[0:1]

	if command != "" {
		os.Args = append(os.Args, command)
	}
	if args != "" {
		os.Args = append(os.Args, args)
	}

	Parse()

	return strings.TrimSpace(buf.String()), exitCode
}

func TestCommandNotFound(t *testing.T) {
	commandName := "instances"
	output, exitCode := ExecParseWithMocks(commandName, "")

	expected := fmt.Sprintf(constants.COMMAND_NOT_FOUND, commandName)
	if output != expected {
		t.Fatalf("\nexpected output:\n\t%s\nactual:\n\t%s", expected, output)
	}
	if exitCode != 1 {
		t.Fatalf("expected exit code of 1, got %d", exitCode)
	}
}

func TestNoCommand(t *testing.T) {
	commandName := ""
	output, exitCode := ExecParseWithMocks(commandName, "")

	trimmedUsage := strings.TrimSpace(usage)
	if output != trimmedUsage {
		t.Fatalf("\nexpected output to be usage\n\texpected:\n\t\t%s\nactual:\n\t\t%s", trimmedUsage, output)
	}
	if exitCode != 1 {
		t.Fatalf("expected exit code of 1, got %d", exitCode)
	}
}

func TestCommandHelp(t *testing.T) {
	commandName := "instances"
	mic := &MockInstancesCmd{
		UsageStr: "<USAGE>",
	}
	RegisterCommand(commandName, mic)

	output, exitCode := ExecParseWithMocks(commandName, "?")

	if output != mic.Usage() {
		t.Fatalf("\nexpected output to be command usage (%s)\nactual:\n\t%s", mic.Usage(), output)
	}

	if exitCode != 0 {
		t.Fatalf("expected exit code of 0, got %d", exitCode)
	}
}

func TestHandlerCalled(t *testing.T) {
	commandName := "instances"
	mic := &MockInstancesCmd{
		UsageStr: `
Usage:
	y2c instances
`,
	}

	RegisterCommand(commandName, mic)

	_, exitCode := ExecParseWithMocks(commandName, "")

	if mic.Calls != 1 {
		t.Fatalf("expected Handler to be called once, was called %d times", mic.Calls)
	}

	if exitCode != 0 {
		t.Fatalf("expected exit code of 0, got %d", exitCode)
	}
}
