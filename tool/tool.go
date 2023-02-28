package tool

import (
	"context"

	"github.com/broothie/ok/argument"
	"github.com/broothie/ok/parameter"
)

type Tool interface {
	Name() string
	CommandName() string
	Filenames() []string
	Extensions() []string
	ProcessFile(path string) ([]Task, error)
}

type Task interface {
	Name() string
	Parameters() parameter.Parameters
	Run(ctx context.Context, args argument.Arguments) error
}
