package ok

import (
	"context"
	"os/exec"

	"github.com/broothie/option"
)

type Task struct {
	Name       string
	RunOptions func(ctx context.Context, args []string) (option.Options[*exec.Cmd], error)
}

type taskInfo struct {
	Task
	tool     *toolInfo
	filePath string
}
