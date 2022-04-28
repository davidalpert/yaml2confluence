package utils

import "github.com/NorthfieldIT/yaml2confluence/internal/resources"

const (
	COMMAND_NOT_FOUND = `"%s" command not found`
	DUPLICATE_TITLE   = `Duplicate title found -- "%s" (%s/) matches "%s" (%s/)`
)

var CHANGE_VERBS = map[resources.ChangeType]string{
	resources.CREATE: "Created",
	resources.UPDATE: "Updated",
	resources.DELETE: "Deleted",
	resources.NOOP:   "Skipped",
}
