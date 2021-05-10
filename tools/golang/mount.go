package golang

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/broothie/now/param"
	"github.com/broothie/now/task"
	"github.com/broothie/now/toolhelp"
)

var funcFinder = regexp.MustCompile(`(?m)^\s*func\s+(\w+)\s*\((.*)\)`)

func (Golang) Mount() ([]task.Task, error) {
	if _, err := os.Open(filename); err != nil {
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

	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, toolhelp.ReadToolFileError{Err: err, Filename: filename}
	}

	fileContents := string(fileBytes)

	matches := funcFinder.FindAllStringSubmatch(fileContents, -1)
	tasks := make([]task.Task, len(matches))
	for i, match := range matches {
		taskName, paramsString := match[1], match[2]

		var paramEntries []string
		if !toolhelp.AllWhitespace(paramsString) {
			paramEntries = strings.Split(paramsString, ",")
		}

		// Loop backwards with a `currentType` because of the whole `func(a, b string)` thing in Go
		params := make([]param.Param, len(paramEntries))
		currentType := ""
		for i := len(paramEntries) - 1; i >= 0; i-- {
			paramEntry := strings.TrimSpace(paramEntries[i])

			chunks := toolhelp.WhitespaceSplitter.Split(paramEntry, 2)
			if len(chunks) > 1 {
				currentType = chunks[1]
			}

			paramName := chunks[0]
			var paramType param.Type
			switch currentType {
			case "interface{}":
				paramType = param.Untyped
			case "bool":
				paramType = param.Bool
			case "int":
				paramType = param.Int
			case "string":
				paramType = param.String
			default:
				return nil, fmt.Errorf("invalid type '%s'", currentType)
			}

			params[i] = param.Param{Name: paramName, Type: paramType}
		}

		tasks[i] = Task{
			Base:         task.NewBaseTask(taskName, filename, ToolName),
			params:       param.Params{PositionalRequired: params},
			fileContents: &fileContents,
		}
	}

	return tasks, nil
}
