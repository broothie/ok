package rake

import (
	"fmt"
	"strings"

	"github.com/broothie/ok/task"
	"github.com/broothie/ok/toolhelp"
)

type Task struct {
	task.Base
	params  task.Parameters
	comment string
}

func (t Task) Comment() string {
	return t.comment
}

func (t Task) Params() task.Parameters {
	return t.params
}

func (t Task) Invoke(args task.Args) task.RunningTask {
	argStrings := make([]string, len(args.Positional))
	for i, arg := range args.Positional {
		argStrings[i] = arg.Value.(string)
	}

	return toolhelp.Exec(ToolName, fmt.Sprintf("%s[%s]", t.Name(), strings.Join(argStrings, ",")))
}
