package ruby

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
	methodFinder = regexp.MustCompile(`(?m)^def\s+(?P<taskName>\w+)\(?(?P<params>.*?)\)?$`)

	positionalMatcher = regexp.MustCompile(`^(?P<paramName>\w+)(?:\s*=\s*(?P<default>.*))?$`)
	keywordMatcher    = regexp.MustCompile(`^(?P<paramName>\w+):(?:\s*(?P<default>.*))?$`)
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

	rawTasks := toolhelp.Scan(file, methodFinder, stringhelp.OctothorpePrefixMatcher)
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

	return tasks, nil
}

func paramListFromParamString(paramsString string) task.Parameters {
	paramStrings := stringhelp.SplitOnCommas(paramsString)
	if len(paramStrings) == 1 && stringhelp.AllWhitespace(paramStrings[0]) {
		return task.Parameters{}
	}

	var params task.Parameters
	for _, paramString := range paramStrings {
		var re *regexp.Regexp
		var isKeyword bool

		if positionalMatcher.MatchString(paramString) {
			re = positionalMatcher
			isKeyword = false
		} else if keywordMatcher.MatchString(paramString) {
			re = keywordMatcher
			isKeyword = true
		} else {
			toolhelp.Warn(ToolName, "error parsing param '%s'", paramString)
			continue
		}

		result := stringhelp.NamedRegexpResult(paramString, re)

		var defaultValue interface{}
		defaultString, defaultPresent := result["default"]
		if defaultPresent && defaultString != "" {
			defaultValue = defaultString
		}

		defaultString = strings.TrimSpace(defaultString)
		if strings.HasPrefix(defaultString, "'") || strings.HasPrefix(defaultString, `"`) {
			defaultValue = strings.Trim(defaultString, `'"`)
		}

		params.ParamList = append(params.ParamList, task.Parameter{
			Name:      result["paramName"],
			Type:      task.Untyped,
			Default:   defaultValue,
			IsKeyword: isKeyword,
		})
	}

	return params
}
