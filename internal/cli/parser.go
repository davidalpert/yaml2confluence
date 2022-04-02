package cli

import (
	"log"
	"os"
	"strings"

	"github.com/NorthfieldIT/yaml2confluence/internal/utils"
	"github.com/docopt/docopt-go"
)

// TODO set this up globally
func init() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
}

var usage = `
y2c

Usage:
	y2c -h | --help
	y2c <command> [? | <args>...]

Options:
	-h --help  		Show this screen.
	<command> ?	  	Usage information for a specific command

Commands:
	render  		Renders a single YAML resource
	instances  		Manage Confluence instance configuration
	
See 'y2c <command> ?' for more information on a specific command.
`

type Command interface {
	Usage() string
	Handler(docopt.Opts)
}

var Commands = map[string]Command{}

func Parse() {
	helpHandlerInvoked := false

	parser := &docopt.Parser{
		OptionsFirst: true,
		HelpHandler: func(err error, _ string) {
			helpHandlerInvoked = true
		},
	}

	args, err := parser.ParseArgs(usage, nil, "")
	if err != nil {
		PrintUsage(usage, 1)
		return
	}
	if helpHandlerInvoked {
		PrintUsage(usage, 0)
		return
	}

	cmd, exists := Commands[args["<command>"].(string)]
	if !exists {
		log.Printf(utils.COMMAND_NOT_FOUND, args["<command>"].(string))
		Exit(1)
		return
	}

	// print command usage if <command> ?
	if args["?"].(bool) {
		PrintUsage(cmd.Usage(), 0)
		return
	}

	// parse command usage and execute handler
	cmdArgs, _ := docopt.ParseDoc(cmd.Usage())
	cmd.Handler(cmdArgs)
}

func PrintUsage(usage string, exitCode int) {
	log.Println(strings.TrimSpace(usage))
	Exit(exitCode)
}

func RegisterCommand(name string, cmd Command) {
	Commands[name] = cmd
}

var exitFunction func(code int) = os.Exit

func Exit(code int) {
	exitFunction(code)
}
