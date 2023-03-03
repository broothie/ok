package tools

import (
	"github.com/broothie/ok/task"
	"github.com/broothie/ok/tool"
)

type Task struct {
	task.Task
	Tool     tool.Tool
	Filename string
}

type Tasks map[string]Task
