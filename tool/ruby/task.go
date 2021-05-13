package ruby

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/broothie/ok/task"
	"github.com/broothie/ok/toolhelp"
)

type Task struct {
	task.Base
	comment string
	params  task.Parameters
}

func (t Task) Comment() string {
	return t.comment
}

func (t Task) Params() task.Parameters {
	return t.params
}

func (t Task) Invoke(args task.Args) *os.Process {
	positionalStrings := make([]string, len(args.Positional))
	for i, arg := range args.Positional {
		positionalStrings[i] = processArg(arg.Value.(string))
	}

	keywordEntries := make([]string, len(args.Keyword))
	counter := 0
	for name, arg := range args.Keyword {
		keywordEntries[counter] = fmt.Sprintf("%s: %s", name, processArg(arg.Value.(string)))
		counter++
	}

	script := fmt.Sprintf("%s(%s)", t.Name(), strings.Join(append(positionalStrings, keywordEntries...), ", "))
	return toolhelp.Exec(ToolName, "-r", fmt.Sprintf("./%s", t.Filename()), "-e", script).Process
}

func processArg(arg string) string {
	if _, err := strconv.ParseFloat(arg, 64); err == nil {
		return arg
	} else if _, err := strconv.Atoi(arg); err == nil {
		return arg
	} else if _, err := strconv.ParseBool(arg); err == nil {
		return arg
	} else {
		return strconv.Quote(arg)
	}
}
