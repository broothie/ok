package ok

import (
	"context"

	"github.com/bobg/errors"
)

type Tool struct {
	Name        string
	CommandName string
	FileGlobs   []string
	ProcessFile func(ctx context.Context, filePath string) ([]Task, error)
}

type toolInfo struct {
	Tool
	commandPath    string
	commandPathErr error
	filePaths      []string
	filePathsErr   error
}

func (i toolInfo) err() error {
	return errors.Join(i.commandPathErr, i.filePathsErr)
}
