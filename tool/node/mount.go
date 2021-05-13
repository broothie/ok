package node

import (
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/broothie/ok/stringhelp"
	"github.com/broothie/ok/task"
	"github.com/broothie/ok/toolhelp"
)

var (
	functionFinder    = regexp.MustCompile(`(?m)^\s*function\s+(?P<taskName>\w+)\s*\((?P<params>.*?)\)\s*{\s*$`)
	positionalMatcher = regexp.MustCompile(`^(?P<paramName>\w+)(?:\s*=\s*(?P<default>.*))?$`)
)

func (t Tool) Mount() ([]task.Task, error) {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, toolhelp.ReadToolFileError{Err: err, Filename: filename}
	}

	if err := t.Check(); err != nil {
		if err == exec.ErrNotFound {
			return nil, toolhelp.CommandNotFoundError{CommandName: ToolName}
		}

		return nil, err
	}

	rawTasks := toolhelp.Scan(file, functionFinder, stringhelp.DoubleSlashPrefixMatcher)
	tasks := make([]task.Task, len(rawTasks))
	for i, rawTask := range rawTasks {
		taskName := rawTask.MatchData["taskName"]
		params := rawTask.MatchData["params"]

		tasks[i] = Task{
			Base:    task.NewBase(taskName, filename, ToolName),
			comment: rawTask.Comment,
			params:  paramListFromParamString(params),
		}
	}

	//fileBytes, err := ioutil.ReadFile(filename)
	//if err != nil {
	//	return nil, toolhelp.ReadToolFileError{Filename: filename, Err: err}
	//}
	//
	//fileContents := string(fileBytes)
	//results := stringhelp.NamedRegexpResults(fileContents, functionFinder)
	//for _, result := range results {
	//	tasks = append(tasks, Task{
	//		Base:         task.NewBase(result["taskName"], filename, ToolName),
	//		params:       paramListFromParamString(result["params"]),
	//		fileContents: &fileContents,
	//	})
	//}

	return tasks, nil
}

func paramListFromParamString(paramsString string) task.Parameters {
	paramStrings := stringhelp.SplitOnCommas(paramsString)
	if len(paramStrings) == 1 && stringhelp.AllWhitespace(paramStrings[0]) {
		return task.Parameters{}
	}

	var params task.Parameters
	for _, paramString := range paramStrings {
		result := stringhelp.NamedRegexpResult(paramString, positionalMatcher)

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
