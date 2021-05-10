package task

import (
	"os"

	"github.com/broothie/now/arg"
	"github.com/broothie/now/param"
)

type Task interface {
	Name() string
	Filename() string
	ToolName() string
	Params() param.Params
	Invoke(args arg.Args) *os.Process
}
