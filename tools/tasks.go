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

func (t Tasks) Task(name string) (Task, bool) {
	task, found := t[name]
	return task, found
}
