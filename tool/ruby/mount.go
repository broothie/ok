package ruby

import (
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/broothie/ok/task"
	"github.com/broothie/ok/tool"
)

var (
	methodFinder  = regexp.MustCompile(`(?m)^\s*def\s+(?P<taskName>\w+)\(?(?P<params>.*?)\)?$`)
	paramSplitter = regexp.MustCompile(`\s*,\s*`)

	positionalMatcher = regexp.MustCompile(`^(?P<paramName>\w+)(?:\s*=\s*(?P<default>.*))?$`)
	keywordMatcher    = regexp.MustCompile(`^(?P<paramName>\w+):(?:\s*(?P<default>.*))?$`)
)

func (t Tool) Mount() ([]task.Task, error) {
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, tool.ReadToolFileError{Err: err, Filename: filename}
	}

	if err := t.Check(); err != nil {
		if err == exec.ErrNotFound {
			return nil, tool.CommandNotFoundError{CommandName: ToolName}
		}

		return nil, err
	}

	var tasks []task.Task
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, tool.ReadToolFileError{Filename: filename, Err: err}
	}

	results := tool.NamedRegexpResults(string(fileBytes), methodFinder)
	for _, result := range results {
		tasks = append(tasks, Task{
			Base:   task.NewBase(result["taskName"], filename, ToolName),
			params: paramListFromParamString(result["params"]),
		})
	}

	return tasks, nil
}

func paramListFromParamString(paramsString string) task.Parameters {
	paramStrings := tool.SplitOnCommas(paramsString)
	if len(paramStrings) == 1 && tool.AllWhitespace(paramStrings[0]) {
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
			tool.Warn(ToolName, "error parsing param '%s'", paramString)
			continue
		}

		result := tool.NamedRegexpResult(paramString, re)

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
