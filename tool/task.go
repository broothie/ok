package tool

import (
	"context"

	"github.com/broothie/ok/argument"
	"github.com/broothie/ok/parameter"
)

type Task interface {
	Name() string
	Parameters() parameter.Parameters
	Run(ctx context.Context, args argument.Arguments) error
}
