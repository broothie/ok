package make

import (
	"os"

	"github.com/broothie/ok/task"
	"github.com/broothie/ok/toolhelp"
)

type Task struct {
	task.Base
	comment string
}

func (t Task) Comment() string {
	return t.comment
}

func (t Task) Invoke(args task.Args) *os.Process {
	return toolhelp.Exec(ToolName, t.Name()).Process
}
