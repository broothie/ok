package bash

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/broothie/ok/task"
	"github.com/broothie/ok/tools/ez"
	"github.com/broothie/ok/util"
)

var taskMatcher = regexp.MustCompile(`^(?P<taskName>\w+)\s*\((?P<params>.*?)\)\s*{\s*$`)

var Bash = ez.Tool{
	ToolName:             "bash",
	CommandName:          "bash",
	ToolFilename:         "Okfile.bash",
	TaskMatcher:          taskMatcher,
	CommentPrefixMatcher: util.OctothorpePrefixMatcher,
	ParamParser: func(paramString string) (task.Parameters, error) {
		paramStrings := util.SplitOnCommas(paramString)
		paramList := make(task.ParamList, len(paramStrings))
		for i, paramString := range paramStrings {
			paramList[i] = task.Parameter{Name: paramString, Type: task.String}
		}

		return paramList.ToParameters(false), nil
	},
	Invoke: func(task ez.Task, args task.Args) (task.RunningTask, error) {
		var argStrings []string
		for _, arg := range args.Positional {
			argStrings = append(argStrings, processArg(arg.Value.(string)))
		}

		script := fmt.Sprintf("%s\n%s %s", *task.FileContents, task.Name(), strings.Join(argStrings, "  "))
		return util.Exec(task.Tool.Name(), "-c", script)
	},
}

func processArg(arg string) string {
	if _, err := strconv.ParseFloat(arg, 64); err == nil {
		return arg
	} else if _, err := strconv.Atoi(arg); err == nil {
		return arg
	} else if b, err := strconv.ParseBool(arg); err == nil {
		if b {
			return "1"
		} else {
			return "0"
		}
	} else {
		return strconv.Quote(arg)
	}
}
