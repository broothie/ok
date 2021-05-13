package yarn

import (
	"os"

	"github.com/broothie/ok/task"
	"github.com/broothie/ok/toolhelp"
)

type Task struct {
	task.Base
}

func (t Task) Comment() string {
	return ""
}

func (Task) Params() task.Parameters {
	return task.Parameters{}
}

func (t Task) Invoke(args task.Args) *os.Process {
	return toolhelp.Exec(ToolName, t.Name()).Process
}
