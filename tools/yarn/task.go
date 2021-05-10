package yarn

import (
	"os"

	"github.com/broothie/okay/arg"
	"github.com/broothie/okay/param"
	"github.com/broothie/okay/task"
	"github.com/broothie/okay/toolhelp"
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
