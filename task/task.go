package task

import "os"

type Task interface {
	Name() string
	Comment() string
	Filename() string
	ToolName() string
	Params() Parameters
	Invoke(args Args) *os.Process
}
