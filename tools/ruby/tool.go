package ruby

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

var definitionRegexp = regexp.MustCompile(`^def (?P<name>\w[a-zA-Z0-9_]*)(?:\(?(?P<paramList>[^)]*)\)?)?$`)

type Tool struct {
	config *tool.Config
}

func New() tool.Tool {
	return Tool{
		config: &tool.Config{
			"extensions": "rb",
			"executable": "ruby",
		},
	}
}

func (Tool) Name() string {
	return "ruby"
}

func (t Tool) Config() *tool.Config {
	return t.config
}

func (t Tool) ProcessFile(path string) ([]task.Task, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read ruby file")
	}

	ruby := string(content)
	lines := strings.Split(ruby, "\n")
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
			switch len(fields) {
			case 1:
				paramName := fields[0]
				params = append(params, task.NewRequired(paramName, task.TypeString))

			case 2:
				paramName, paramDefault := fields[0], fields[1]
				params = append(params, task.NewOptional(strings.TrimSuffix(paramName, ":"), parseType(paramDefault), paramDefault))

			default:
				return nil, fmt.Errorf("invalid parameter %q", param)
			}
		}

		description := ""
		if i != 0 && strings.HasPrefix(lines[i-1], "#") {
			description = lines[i-1]
		}

		tasks = append(tasks, Task{
			Tool:        t,
			name:        name,
			description: description,
			parameters:  params,
			filename:    path,
		})
	}

	return tasks, nil
}

func parseType(param string) task.Type {
	if param == "false" || param == "true" {
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
