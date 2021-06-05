package tool

import "github.com/broothie/ok/task"

type Tool interface {
	Name() string
	Init() error
	Check() error
	Config() interface{}
	Mount() ([]task.Task, error)
}
