package util

import (
	"fmt"
	"strings"

	"github.com/thoas/go-funk"
)

type CommandNotFoundError struct {
	CommandName string
}

func (e CommandNotFoundError) Error() string {
	return fmt.Sprintf("command %q not found", e.CommandName)
}

type ReadToolFileError struct {
	Filename string
	Err      error
}

func (e ReadToolFileError) Error() string {
	return fmt.Sprintf("error reading file %q: %v", e.Filename, e.Err)
}

type ErrorGroup []error

func (eg *ErrorGroup) Error() string {
	return fmt.Sprintf("multiple errors occurred:\n%s", strings.Join(funk.Map(eg, func(err error) string { return err.Error() }).([]string), "\n"))
}

func (eg *ErrorGroup) Add(errs ...error) {
	*eg = append(*eg, errs...)
}

func (eg *ErrorGroup) NilIfEmpty() error {
	if len(*eg) == 0 {
		return nil
	}

	return eg
}
