package yarn

import (
	"os"

	"github.com/broothie/now/arg"
	"github.com/broothie/now/param"
	"github.com/broothie/now/task"
	"github.com/broothie/now/toolhelp"
)

type Task struct {
	task.Base
}

func (Task) Params() param.Params {
	return param.Params{}
}

func (t Task) Invoke(args arg.Args) *os.Process {
	return toolhelp.Exec(ToolName, t.Name()).Process
}
