package okay

import "github.com/broothie/okay/task"

type Tool interface {
	Name() string
	Init() error
	Check() error
	Mount() ([]task.Task, error)
}
