package rake

import (
	"fmt"
	"strings"

	"github.com/broothie/ok/task"
	"github.com/broothie/ok/util"
)

type Task struct {
	task.Base
	params  task.Parameters
	comment string
	tool    *Tool
}

func (t Task) Comment() string {
	return t.comment
}

func (t Task) Params() task.Parameters {
	return t.params
}

func (t Task) Invoke(args task.Args) (task.RunningTask, error) {
	argStrings := make([]string, len(args.Positional))
	for i, arg := range args.Positional {
		argStrings[i] = arg.Value.(string)
	}

	command := "rake"
	taskString := fmt.Sprintf("%s[%s]", t.Name(), strings.Join(argStrings, ","))
	rest := []string{taskString}
	if t.tool.Config.Bundler != nil && *t.tool.Config.Bundler {
		command = "bundle"
		rest = []string{"exec", "rake", taskString}
	}

	return util.Exec(command, rest...)
}
