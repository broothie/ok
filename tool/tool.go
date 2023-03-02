package tool

import "github.com/broothie/ok/task"

type Tool interface {
	Name() string
	Executable() string
	Filenames() []string
	Extensions() []string
	ProcessFile(path string) ([]task.Task, error)
}
