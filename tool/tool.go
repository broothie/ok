package tool

import "github.com/broothie/okay/task"

type Tool interface {
	Init() error
	Mount() ([]task.Task, error)
}
