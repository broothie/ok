package dockercompose

import (
	"github.com/broothie/ok/task"
	"github.com/broothie/ok/util"
)

type Task struct {
	task.Base
}

func (Task) Comment() string {
	return ""
}

func (Task) Params() task.Parameters {
	return task.Parameters{Forward: true}
}

func (t Task) Invoke(args task.Args) (task.RunningTask, error) {
	argStrings := []string{"run", t.Name()}
	return util.Exec(ToolName, append(argStrings, args.Forwards...)...)
}
