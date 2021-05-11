package yarn

import (
	"os"

	"github.com/broothie/okay/tool"

	"github.com/broothie/okay/task"
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
