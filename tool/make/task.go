package make

import (
	"os"

	"github.com/broothie/now/arg"
	"github.com/broothie/now/param"
	"github.com/broothie/now/task"
	"github.com/broothie/now/tool"
)

type Task struct {
	task.Base
}

func (Task) Params() param.Params {
	return param.Params{}
}

func (t Task) Invoke(args arg.Task) *os.Process {
	return tool.Exec(ToolName, t.Name()).Process
}
