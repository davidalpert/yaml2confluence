package main

import (
	"github.com/NorthfieldIT/yaml2confluence/internal/cli"
	_ "github.com/NorthfieldIT/yaml2confluence/internal/commands"
)

func main() {
	cli.Parse()
}
