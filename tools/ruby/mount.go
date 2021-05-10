package ruby

import (
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/broothie/now/param"
	"github.com/broothie/now/task"
	"github.com/broothie/now/toolhelp"
)

var (
	methodFinder  = regexp.MustCompile(`(?m)^\s*def\s+(?P<taskName>\w+)\(?(?P<params>.*?)\)?$`)
	paramSplitter = regexp.MustCompile(`\s*,\s*`)

	positionalMatcher = regexp.MustCompile(`^(?P<paramName>\w+)(?:\s*=\s*(?P<default>.*))?$`)
	keywordMatcher    = regexp.MustCompile(`^(?P<paramName>\w+):(?:\s*(?P<default>.*))?$`)
)

func (Ruby) Mount() ([]task.Task, error) {
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, toolhelp.ReadToolFileError{Err: err, Filename: filename}
	}

	if _, err := exec.LookPath(ToolName); err != nil {
		if err == exec.ErrNotFound {
			return nil, toolhelp.CommandNotFoundError{CommandName: ToolName}
		}

		return nil, err
	}

	var tasks []task.Task
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, toolhelp.ReadToolFileError{Filename: filename, Err: err}
	}

	results := toolhelp.NamedRegexpResults(string(fileBytes), methodFinder)
	for _, result := range results {
		tasks = append(tasks, Task{
			Base:   task.NewBaseTask(result["taskName"], filename, ToolName),
			params: paramListFromParamString(result["params"]),
		})
	}

	return tasks, nil
}

func paramListFromParamString(paramsString string) param.Params {
	paramStrings := paramSplitter.Split(paramsString, -1)
	if len(paramStrings) == 1 && toolhelp.AllWhitespace(paramStrings[0]) {
		paramStrings = nil
	}

	var params param.Params
	for _, paramString := range paramStrings {
		var re *regexp.Regexp
		var paramListWithoutDefault *[]param.Param
		var paramListWithDefault *[]param.Param

		if positionalMatcher.MatchString(paramString) {
			re = positionalMatcher
			paramListWithoutDefault = &params.PositionalRequired
			paramListWithDefault = &params.PositionalOptional
		} else if keywordMatcher.MatchString(paramString) {
			re = keywordMatcher
			paramListWithoutDefault = &params.KeywordRequired
			paramListWithDefault = &params.KeywordOptional
		} else {
			toolhelp.Warn(ToolName, "error parsing param '%s'", paramString)
			continue
		}

		var paramList *[]param.Param
		var defaultString string
		var defaultExists bool
		result := toolhelp.NamedRegexpResult(paramString, re)
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

		*paramList = append(*paramList, param.Param{
			Name:    result["paramName"],
			Type:    param.Untyped,
			Default: defaultValue,
		})
	}

	return params
}
