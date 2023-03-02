package task

import "context"

type Task interface {
	Name() string
	Parameters() Parameters
	Run(ctx context.Context, args Arguments) error
}
