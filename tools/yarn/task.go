package yarn

import (
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

func (t Task) Invoke(args task.Args) task.RunningTask {
	return toolhelp.Exec(ToolName, t.Name())
}
