package golang

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/broothie/ok/task"
	"github.com/broothie/ok/tool"
	"github.com/broothie/ok/util"
	"github.com/pkg/errors"
)

const commentPrefix = "//"

var definitionRegexp = regexp.MustCompile(`^func (?P<name>\w[a-zA-Z0-9_]*)\((?P<paramList>[^)]*)\) \{$`)

type Tool struct {
	config *tool.Config
}

func New() tool.Tool {
	return Tool{
		config: &tool.Config{
			"extensions": "go",
			"executable": "go",
		},
	}
}

func (Tool) Name() string {
	return "go"
}

func (t Tool) Config() *tool.Config {
	return t.config
}

func (t Tool) ProcessFile(path string) ([]task.Task, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read ruby file")
	}

	goCode := string(content)
	lines := strings.Split(goCode, "\n")
	var tasks []task.Task
	for i, line := range lines {
		captures := util.NamedCaptureGroups(definitionRegexp, line)
		if len(captures) == 0 {
			continue
		}

		name := captures["name"]
		paramList := captures["paramList"]
		var params task.Parameters
		for _, param := range util.SplitCommaList(paramList) {
			fields := strings.Fields(param)
			paramName, paramType := fields[0], fields[1]
			typ, err := parseType(paramType)
			if err != nil {
				return nil, err
			}

			params = append(params, task.NewRequired(paramName, typ))
		}

		description := ""
		if i != 0 {
			description = util.ExtractComment(lines[i-1], "//")
		}

		tasks = append(tasks, Task{
			Tool:        t,
			name:        name,
			description: description,
			parameters:  params,
			filename:    path,
			goCode:      &goCode,
		})
	}

	return tasks, nil
}

func parseType(paramType string) (task.Type, error) {
	var typ task.Type
	switch paramType {
	case "bool":
		typ = task.TypeBool
	case "float32", "float64":
		typ = task.TypeFloat
	case "int", "int8", "int16", "int32", "int64":
		typ = task.TypeInt
	case "string":
		typ = task.TypeString
	default:
		return "", fmt.Errorf("invalid type for: %s", paramType)
	}

	return typ, nil
}
