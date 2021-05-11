package golang

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/broothie/ok/task"
	"github.com/broothie/ok/tool"
)

var funcFinder = regexp.MustCompile(`(?m)^\s*func\s+(\w+)\s*\((.*)\)`)

func (t Tool) Mount() ([]task.Task, error) {
	if _, err := os.Open(filename); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, tool.ReadToolFileError{Err: err, Filename: filename}
	}

	if err := t.Check(); err != nil {
		return nil, err
	}

	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, tool.ReadToolFileError{Err: err, Filename: filename}
	}

	fileContents := string(fileBytes)

	matches := funcFinder.FindAllStringSubmatch(fileContents, -1)
	tasks := make([]task.Task, len(matches))
	for i, match := range matches {
		taskName, paramsString := match[1], match[2]

		var paramEntries []string
		if !tool.AllWhitespace(paramsString) {
			paramEntries = strings.Split(paramsString, ",")
		}

		// Loop backwards with a `currentType` because of the whole `func(a, b string)` thing in Go
		params := make([]task.Parameter, len(paramEntries))
		currentType := ""
		for i := len(paramEntries) - 1; i >= 0; i-- {
			paramEntry := strings.TrimSpace(paramEntries[i])

			chunks := tool.WhitespaceSplitter.Split(paramEntry, 2)
			if len(chunks) > 1 {
				currentType = chunks[1]
			}

			paramName := chunks[0]
			var paramType task.Type
			switch currentType {
			case "interface{}":
				paramType = task.Untyped
			case "bool":
				paramType = task.Bool
			case "int":
				paramType = task.Int
			case "float64", "float32":
				paramType = task.Float
			case "string":
				paramType = task.String
			default:
				return nil, fmt.Errorf("invalid type '%s'", currentType)
			}

			params[i] = task.Parameter{Name: paramName, Type: paramType}
		}

		tasks[i] = Task{
			Base:         task.NewBase(taskName, filename, ToolName),
			params:       task.Parameters{PositionalRequired: params},
			fileContents: &fileContents,
		}
	}

	return tasks, nil
}
