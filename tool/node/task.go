package node

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/broothie/ok/task"
	"github.com/broothie/ok/tool"
)

type Task struct {
	task.Base
	params       task.Parameters
	fileContents *string
}

func (t Task) Params() task.Parameters {
	return t.params
}

func (t Task) Invoke(args task.Args) *os.Process {
	positionalStrings := make([]string, len(args.Positional))
	for i, arg := range args.Positional {
		positionalStrings[i] = processArg(arg.Value.(string))
	}

	argString := strings.Join(positionalStrings, ", ")
	script := fmt.Sprintf("%s; %s(%s)", *t.fileContents, t.Name(), argString)
	return tool.Exec(ToolName, "-e", script).Process
}

func processArg(arg string) string {
	if _, err := strconv.ParseFloat(arg, 64); err == nil {
		return arg
	} else if _, err := strconv.Atoi(arg); err == nil {
		return arg
	} else if _, err := strconv.ParseBool(arg); err == nil {
		return arg
	} else {
		return fmt.Sprintf(`"%s"`, arg)
	}
}
