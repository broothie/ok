package ruby

import (
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/broothie/now/param"
	"github.com/broothie/now/task"
	"github.com/broothie/now/tool"
)

const ToolName = "ruby"

var (
	methodFinder  = regexp.MustCompile(`(?m)^\s*def\s+(?P<taskName>\w+)\(?(?P<params>.*?)\)?$`)
	paramSplitter = regexp.MustCompile(`\s*,\s*`)

	positionalMatcher = regexp.MustCompile(`^(?P<paramName>\w+)(?:\s*=\s*(?P<default>.*))?$`)
	keywordMatcher    = regexp.MustCompile(`^(?P<paramName>\w+):(?:\s*(?P<default>.*))?$`)
)

func Mount() ([]task.Task, error) {
	filenames, _ := filepath.Glob("Nowfile.rb")
	if len(filenames) == 0 {
		return nil, nil
	}

	if _, err := exec.LookPath(ToolName); err != nil {
		if err == exec.ErrNotFound {
			return nil, tool.CommandNotFoundError{CommandName: ToolName}
		}

		return nil, err
	}

	var tasks []task.Task
	for _, filename := range filenames {
		fileBytes, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, tool.ReadToolFileError{Filename: filename, Err: err}
		}

		results := tool.NamedRegexpResults(string(fileBytes), methodFinder)
		for _, result := range results {
			tasks = append(tasks, Task{
				Base:   task.NewBaseTask(result["taskName"], filename, ToolName),
				params: paramListFromParamString(result["params"]),
			})
		}
	}

	return tasks, nil
}

func paramListFromParamString(paramsString string) param.Params {
	paramStrings := paramSplitter.Split(paramsString, -1)

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
			tool.Warn(ToolName, "error parsing param '%s'", paramString)
			continue
		}

		var paramList *[]param.Param
		var defaultString string
		var defaultExists bool
		result := tool.NamedRegexpResult(paramString, re)
		if defaultString, defaultExists = result["default"]; defaultExists && defaultString != "" {
			paramList = paramListWithDefault
		} else {
			paramList = paramListWithoutDefault
		}

		defaultString = strings.TrimSpace(defaultString)

		var defaultValue interface{}
		if defaultString != "" {
			if strings.HasPrefix(defaultString, `"`) || strings.HasPrefix(defaultString, "'") {
				defaultValue = strings.Trim(defaultString, `'"`)
			} else if defaultString == "true" {
				defaultValue = true
			} else if defaultString == "false" {
				defaultValue = false
			} else if f, err := strconv.ParseFloat(defaultString, 64); err == nil {
				defaultValue = f
			} else if i, err := strconv.Atoi(defaultString); err == nil {
				defaultValue = i
			}
		}

		*paramList = append(*paramList, param.Param{
			Name:    result["paramName"],
			Type:    param.Untyped,
			Default: defaultValue,
		})
	}

	return params
}
