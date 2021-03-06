package golang

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/broothie/ok/task"
	"github.com/broothie/ok/util"
)

var funcFinder = regexp.MustCompile(`(?m)^\s*func\s+(?P<taskName>\w+)\s*\((?P<params>.*?)\)`)

func (t Tool) Mount() ([]task.Task, error) {
	if _, err := os.Open(filename); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, util.ReadToolFileError{Err: err, Filename: filename}
	}

	if err := t.Check(); err != nil {
		return nil, err
	}

	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, util.ReadToolFileError{Err: err, Filename: filename}
	}

	fileContents := string(fileBytes)
	rawTasks := util.Scan(bytes.NewBuffer(fileBytes), funcFinder, util.DoubleSlashPrefixMatcher)
	tasks := make([]task.Task, len(rawTasks))
	for i, rawTask := range rawTasks {
		taskName := rawTask.MatchData["taskName"]
		paramsString := rawTask.MatchData["params"]

		var paramEntries []string
		if !util.AllWhitespace(paramsString) {
			paramEntries = util.SplitOnCommas(paramsString)
		}

		// Loop backwards with a `currentType` because of the whole `func(a, b string)` thing in Go
		params := make(task.ParamList, len(paramEntries))
		currentType := ""
		for i := len(paramEntries) - 1; i >= 0; i-- {
			paramEntry := strings.TrimSpace(paramEntries[i])

			chunks := util.Whitespace.Split(paramEntry, 2)
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
				return nil, fmt.Errorf("invalid type %q", currentType)
			}

			params[i] = task.Parameter{Name: paramName, Type: paramType}
		}

		tasks[i] = Task{
			Base:         task.NewBase(taskName, filename, ToolName),
			params:       params.ToParameters(false),
			comment:      rawTask.Comment,
			fileContents: &fileContents,
		}
	}

	return tasks, nil
}
