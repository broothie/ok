package tool

import (
	"context"
	"os/exec"

	"github.com/broothie/option"
)

type Tool struct {
	Name        string
	CommandName string
	FileGlobs   []string
	ParseFile   func(ctx context.Context, filePath string) ([]Task, error)
}

type Task struct {
	Name       string
	RunOptions func(ctx context.Context, args []string) (option.Options[*exec.Cmd], error)
}
