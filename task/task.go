package task

import (
	"os"

	"github.com/broothie/okay/arg"
	"github.com/broothie/okay/param"
)

type Task interface {
	Name() string
	Filename() string
	ToolName() string
	Params() param.Params
	Invoke(args arg.Args) *os.Process
}
