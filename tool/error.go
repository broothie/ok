package tool

import "fmt"

type CommandNotFoundError struct {
	CommandName string
}

func (e CommandNotFoundError) Error() string {
	return fmt.Sprintf("command '%s' not found", e.CommandName)
}

type ReadToolFileError struct {
	Err      error
	Filename string
}

func (e ReadToolFileError) Error() string {
	return fmt.Sprintf("error reading file '%s': %v", e.Filename, e.Err)
}
