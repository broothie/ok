package dockercompose

import (
	"os"

	"github.com/broothie/ok/task"
	"github.com/broothie/ok/tool"
)

type Task struct {
	task.Base
}

func (Task) Params() task.Parameters {
	return task.Parameters{Forward: true}
}

func (t Task) Invoke(args task.Args) *os.Process {
	argStrings := []string{"run", t.Name()}
	return tool.Exec(ToolName, append(argStrings, args.Forwards...)...).Process
}
