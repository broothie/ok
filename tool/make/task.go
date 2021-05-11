package make

import (
	"os"

	"github.com/broothie/ok/task"
	"github.com/broothie/ok/tool"
)

type Task struct {
	task.Base
}

func (Task) Params() task.Parameters {
	return task.Parameters{}
}

func (t Task) Invoke(args task.Args) *os.Process {
	return tool.Exec(ToolName, t.Name()).Process
}
