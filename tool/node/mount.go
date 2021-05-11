package node

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
	functionFinder = regexp.MustCompile(`(?m)^\s*function\s+(?P<taskName>\w+)\s*\((?P<params>.*?)\)\s*{\s*$`)
	paramSplitter  = regexp.MustCompile(`\s*,\s*`)

	positionalMatcher = regexp.MustCompile(`^(?P<paramName>\w+)(?:\s*=\s*(?P<default>.*))?$`)
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

	fileContents := string(fileBytes)
	results := tool.NamedRegexpResults(fileContents, functionFinder)
	for _, result := range results {
		tasks = append(tasks, Task{
			Base:         task.NewBase(result["taskName"], filename, ToolName),
			params:       paramListFromParamString(result["params"]),
			fileContents: &fileContents,
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
		var paramList *[]task.Parameter
		var defaultString string
		var defaultExists bool
		result := tool.NamedRegexpResult(paramString, positionalMatcher)
		if defaultString, defaultExists = result["default"]; defaultExists && defaultString != "" {
			paramList = &params.PositionalOptional
		} else {
			paramList = &params.PositionalRequired
		}

		defaultString = strings.TrimSpace(defaultString)
		var defaultValue interface{} = defaultString
		if strings.HasPrefix(defaultString, "'") || strings.HasPrefix(defaultString, `"`) || strings.HasPrefix(defaultString, "`") {
			defaultValue = strings.Trim(defaultString, "`'\"")
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
