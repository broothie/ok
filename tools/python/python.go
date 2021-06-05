package python

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/broothie/ok/task"
	"github.com/broothie/ok/tools/ez"
	"github.com/broothie/ok/util"
)

var (
	taskMatcher  = regexp.MustCompile(`^def\s+(?P<taskName>\w+)\s*\((?P<params>.*?)\)\s*:\s*$`)
	paramMatcher = regexp.MustCompile(`^(?P<paramName>\w+)(?:\s*=\s*(?P<default>.*))?$`)
)

var Python = ez.Tool{
	ToolName:             "python",
	CommandName:          "python",
	ToolFilename:         "Okfile.py",
	TaskMatcher:          taskMatcher,
	CommentPrefixMatcher: util.OctothorpePrefixMatcher,
	ParamParser: func(paramString string) (task.Parameters, error) {
		paramStrings := util.SplitOnCommas(paramString)
		paramList := make(task.ParamList, len(paramStrings))
		for i, paramString := range paramStrings {
			result := util.NamedRegexpResult(paramString, paramMatcher)
			paramName := result["paramName"]

			var defaultValue interface{}
			defaultString, defaultExists := result["default"]
			if defaultExists && defaultString != "" {
				defaultValue = defaultString
			}

			paramList[i] = task.Parameter{
				Name:      paramName,
				Type:      task.Untyped,
				Default:   defaultValue,
				IsKeyword: defaultString != "",
			}
		}

		return paramList.ToParameters(false), nil
	},
	Invoke: func(task ez.Task, args task.Args) task.RunningTask {
		var argStrings []string
		for _, arg := range args.Positional {
			argStrings = append(argStrings, processArg(arg.Value.(string)))
		}

		for name, arg := range args.Keyword {
			argStrings = append(argStrings, fmt.Sprintf("%s=%s", name, processArg(arg.Value.(string))))
		}

		script := fmt.Sprintf("%s\n%s(%s)", *task.FileContents, task.Name(), strings.Join(argStrings, ", "))
		return util.Exec(task.ToolName(), "-c", script)
	},
}

func processArg(arg string) string {
	if _, err := strconv.ParseFloat(arg, 64); err == nil {
		return arg
	} else if _, err := strconv.Atoi(arg); err == nil {
		return arg
	} else if b, err := strconv.ParseBool(arg); err == nil {
		if b {
			return "True"
		} else {
			return "False"
		}
	} else {
		return strconv.Quote(arg)
	}
}
