package make

import (
	"context"

	"github.com/broothie/ok/argument"
	"github.com/broothie/ok/parameter"
	"github.com/broothie/ok/util"
	"github.com/pkg/errors"
)

type Task struct {
	name string
}

func (r Task) Name() string {
	return r.name
}

func (r Task) Parameters() parameter.Parameters {
	return nil
}

func (r Task) Run(ctx context.Context, _ argument.Arguments) error {
	if err := util.CommandContext(ctx, "make", r.name).Run(); err != nil {
		return errors.Wrap(err, "failed to run make rule")
	}

	return nil
}
