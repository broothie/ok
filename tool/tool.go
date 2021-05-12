package tool

import "github.com/broothie/ok/task"

type Tool interface {
	Name() string
	Init() error
	Check() error
	Mount() ([]task.Task, error)
}
