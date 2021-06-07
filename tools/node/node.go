package node

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
	functionFinder    = regexp.MustCompile(`(?m)^\s*function\s+(?P<taskName>\w+)\s*\((?P<params>.*?)\)\s*{\s*$`)
	positionalMatcher = regexp.MustCompile(`^(?P<paramName>\w+)(?:\s*=\s*(?P<default>.*))?$`)

	Node = ez.Tool{
		ToolName:             "node",
		CommandName:          "node",
		ToolFilename:         "Okfile.js",
		TaskMatcher:          functionFinder,
		CommentPrefixMatcher: util.DoubleSlashPrefixMatcher,
		ParamParser: func(paramString string) (task.Parameters, error) {
			return paramListFromParamString(paramString), nil
		},
		Invoke: func(task ez.Task, args task.Args) (task.RunningTask, error) {
			positionalStrings := make([]string, len(args.Positional))
			for i, arg := range args.Positional {
				positionalStrings[i] = processArg(arg.Value.(string))
			}

			argString := strings.Join(positionalStrings, ", ")
			script := fmt.Sprintf("%s; %s(%s)", *task.FileContents, task.Name(), argString)
			return util.Exec("node", "-e", script)
		},
	}
)

func paramListFromParamString(paramsString string) task.Parameters {
	paramStrings := util.SplitOnCommas(paramsString)
	if len(paramStrings) == 1 && util.AllWhitespace(paramStrings[0]) {
		return task.Parameters{}
	}

	var params task.Parameters
	for _, paramString := range paramStrings {
		result := util.NamedRegexpResult(paramString, positionalMatcher)

		var defaultValue interface{}
		defaultString, defaultPresent := result["default"]
		if defaultPresent && defaultString != "" {
			defaultValue = defaultString
		}

		if strings.HasPrefix(defaultString, "'") || strings.HasPrefix(defaultString, `"`) || strings.HasPrefix(defaultString, "`") {
			defaultValue = strings.Trim(defaultString, "`'\"")
		}

		params.ParamList = append(params.ParamList, task.Parameter{
			Name:    result["paramName"],
			Type:    task.Untyped,
			Default: defaultValue,
		})
	}

	return params
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
