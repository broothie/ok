package yarn

import (
	"github.com/broothie/ok/task"
	"github.com/broothie/ok/util"
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

func (t Task) Invoke(task.Args) (task.RunningTask, error) {
	return util.Exec(ToolName, t.Name())
}
