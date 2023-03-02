package golang

import (
	"os"
	"regexp"
	"strings"

	"github.com/broothie/ok/logger"
	"github.com/broothie/ok/task"
	"github.com/broothie/ok/util"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

var definitionRegexp = regexp.MustCompile(`^func (?P<name>\w[a-zA-Z0-9_]*)\((?P<paramList>[^)]*)\) \{$`)

type Tool struct{}

func (Tool) Name() string {
	return "Go"
}

func (Tool) Executable() string {
	return "go"
}

func (Tool) Filenames() []string {
	return nil
}

func (Tool) Extensions() []string {
	return []string{"go"}
}

func (Tool) ProcessFile(path string) ([]task.Task, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read ruby file")
	}

	goCode := string(content)

	return lo.FilterMap(strings.Split(goCode, "\n"), func(line string, _ int) (task.Task, bool) {
		captures := util.NamedCaptureGroups(definitionRegexp, line)
		if len(captures) == 0 {
			return nil, false
		}

		taskName := captures["name"]
		paramList := captures["paramList"]
		var params task.Parameters
		for _, param := range util.SplitCommaParamList(paramList) {
			fields := strings.Fields(param)
			paramName, paramType := fields[0], fields[1]
			params = append(params, task.NewRequired(paramName, parseType(paramType)))
		}

		return Task{name: taskName, parameters: params, filename: path, goCode: &goCode}, true
	}), nil
}

func parseType(paramType string) task.Type {
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
		logger.Log.Println("invalid type for: %s", paramType)
	}

	return typ
}
