package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/docopt/docopt-go"
)

var usage = `
y2c

Usage:
	y2c -h | --help
	y2c <command> [? | <args>...]

Options:
	-h --help  		Show this screen.
	<command> ?  	Usage information for a specific command

Commands:
	render  		Renders a single YAML resource
	help  			Displays the usage for a command
	instances  		Manage Confluence instance configuration
	
See 'y2c <command> ?' for more information on a specific command.
`

type Command struct {
	usage   string
	handler func(docopt.Opts)
}

var commands = map[string]Command{}

func Parse() {
	parser := &docopt.Parser{OptionsFirst: true}
	args, _ := parser.ParseArgs(usage, nil, "")

	// print usage if no command provided
	cmdName, exists := args["<command>"].(string)
	if !exists {
		PrintUsage(usage, 1)
	}

	// print usage if command doesn't exist
	cmd, exists := commands[cmdName]
	if !exists {
		PrintUsage(usage, 1)
	}

	// print command usage if <command> ?
	if args["?"].(bool) {
		PrintUsage(cmd.usage, 0)
	}

	// parse command usage and execute handler
	cmdArgs, _ := docopt.ParseDoc(cmd.usage)
	cmd.handler(cmdArgs)
}

func PrintUsage(usage string, exitCode int) {
	fmt.Println(strings.TrimSpace(usage))
	os.Exit(exitCode)
}

func RegisterCommand(name string, usage string, handler func(docopt.Opts)) {
	commands[name] = Command{usage, handler}
}

func GetCommandUsage(name string) string {
	return commands[name].usage
}
