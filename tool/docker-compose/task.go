package dockercompose

import (
	"os"

	"github.com/broothie/ok/task"
	"github.com/broothie/ok/toolhelp"
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

func (t Task) Invoke(args task.Args) *os.Process {
	argStrings := []string{"run", t.Name()}
	return toolhelp.Exec(ToolName, append(argStrings, args.Forwards...)...).Process
}
