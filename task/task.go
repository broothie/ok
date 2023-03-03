package task

import "context"

type Task interface {
	Name() string
	Parameters() Parameters
	Description() string
	Run(ctx context.Context, args Arguments) error
}
