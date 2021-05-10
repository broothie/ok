package tool

import "github.com/broothie/now/task"

type Tool interface {
	Init() error
	Mount() ([]task.Task, error)
}
