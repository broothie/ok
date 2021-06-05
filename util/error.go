package util

import "fmt"

type CommandNotFoundError struct {
	CommandName string
}

func (e CommandNotFoundError) Error() string {
	return fmt.Sprintf("command '%s' not found", e.CommandName)
}

type ReadToolFileError struct {
	Filename string
	Err      error
}

func (e ReadToolFileError) Error() string {
	return fmt.Sprintf("error reading file '%s': %v", e.Filename, e.Err)
}
