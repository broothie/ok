package ruby

import (
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/broothie/okay/task"
	"github.com/broothie/okay/tool"
)

var (
	methodFinder  = regexp.MustCompile(`(?m)^\s*def\s+(?P<taskName>\w+)\(?(?P<params>.*?)\)?$`)
	paramSplitter = regexp.MustCompile(`\s*,\s*`)

	positionalMatcher = regexp.MustCompile(`^(?P<paramName>\w+)(?:\s*=\s*(?P<default>.*))?$`)
	keywordMatcher    = regexp.MustCompile(`^(?P<paramName>\w+):(?:\s*(?P<default>.*))?$`)
)

func (r Ruby) Mount() ([]task.Task, error) {
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, tool.ReadToolFileError{Err: err, Filename: filename}
	}

	if err := r.Check(); err != nil {
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
	paramStrings := paramSplitter.Split(paramsString, -1)
	if len(paramStrings) == 1 && tool.AllWhitespace(paramStrings[0]) {
		paramStrings = nil
	}

	var params task.Parameters
	for _, paramString := range paramStrings {
		var re *regexp.Regexp
		var paramListWithoutDefault *[]task.Parameter
		var paramListWithDefault *[]task.Parameter

		if positionalMatcher.MatchString(paramString) {
			re = positionalMatcher
			paramListWithoutDefault = &params.PositionalRequired
			paramListWithDefault = &params.PositionalOptional
		} else if keywordMatcher.MatchString(paramString) {
			re = keywordMatcher
			paramListWithoutDefault = &params.KeywordRequired
			paramListWithDefault = &params.KeywordOptional
		} else {
			tool.Warn(ToolName, "error parsing param '%s'", paramString)
			continue
		}

		var paramList *[]task.Parameter
		var defaultString string
		var defaultExists bool
		result := tool.NamedRegexpResult(paramString, re)
		if defaultString, defaultExists = result["default"]; defaultExists && defaultString != "" {
			paramList = paramListWithDefault
		} else {
			paramList = paramListWithoutDefault
		}

		defaultString = strings.TrimSpace(defaultString)
		var defaultValue interface{} = defaultString
		if strings.HasPrefix(defaultString, "'") || strings.HasPrefix(defaultString, `"`) {
			defaultValue = strings.Trim(defaultString, `'"`)
		} else if defaultString == "" {
			defaultValue = nil
		}

		*paramList = append(*paramList, task.Parameter{
			Name:    result["paramName"],
			Type:    task.Untyped,
			Default: defaultValue,
		})
	}

	return params
}
