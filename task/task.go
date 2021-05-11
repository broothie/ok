package task

import (
	"os"
)

type Task interface {
	Name() string
	Filename() string
	ToolName() string
	Params() Parameters
	Invoke(args Args) *os.Process
}
