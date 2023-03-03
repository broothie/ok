package python

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/broothie/ok/task"
	"github.com/broothie/ok/tool"
	"github.com/broothie/ok/util"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

var definitionRegexp = regexp.MustCompile(`^def (?P<name>\w[a-zA-Z0-9_]*)\((?P<paramList>[^)]*)\):$`)

type Tool struct {
	config *tool.Config
}

func New() tool.Tool {
	return Tool{
		config: &tool.Config{
			"extensions": "py",
			"executable": "python",
		},
	}
}

func (Tool) Name() string {
	return "python"
}

func (t Tool) Config() *tool.Config {
	return t.config
}

func (t Tool) ProcessFile(path string) ([]task.Task, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read python file")
	}

	python := string(content)
	var tasks []task.Task
	for _, line := range strings.Split(python, "\n") {
		captures := util.NamedCaptureGroups(definitionRegexp, line)
		if len(captures) == 0 {
			continue
		}

		name := captures["name"]
		paramList := captures["paramList"]
		var params task.Parameters
		for _, param := range util.SplitCommaList(paramList) {
			fields := strings.SplitN(param, "=", 2)
			switch len(fields) {
			case 1:
				paramName := strings.TrimSpace(fields[0])
				params = append(params, task.NewRequired(paramName, task.TypeString))

			case 2:
				paramName, paramDefault := strings.TrimSpace(fields[0]), strings.TrimSpace(fields[1])
				params = append(params, task.NewOptional(paramName, parseType(paramDefault), paramDefault))

			default:
				return nil, fmt.Errorf("invalid parameter %q", param)
			}
		}

		tasks = append(tasks, Task{
			Tool:       t,
			name:       name,
			parameters: params,
			pythonCode: &python,
		})
	}

	return tasks, nil
}

func parseType(param string) task.Type {
	if param == "False" || param == "True" {
		return task.TypeBool
	} else if lo.Every([]rune("1234567890_"), []rune(param)) {
		return task.TypeInt
	} else if lo.Every([]rune("1234567890_."), []rune(param)) {
		return task.TypeFloat
	} else if (strings.HasPrefix(param, `"`) && strings.HasSuffix(param, `"`)) || (strings.HasPrefix(param, `'`) && strings.HasSuffix(param, `'`)) {
		return task.TypeString
	} else {
		panic(fmt.Sprintf("invalid type %q", param))
	}
}
