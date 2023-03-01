package python

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/broothie/ok/logger"
	"github.com/broothie/ok/parameter"
	"github.com/broothie/ok/tool"
	"github.com/broothie/ok/util"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

var (
	definitionRegexp  = regexp.MustCompile(`^def (?P<name>\w[a-zA-Z0-9_]*)\((?P<paramList>[^)]*)\):$`)
	equalSignSplitter = regexp.MustCompile(`\s*=\s*`)
)

type Tool struct{}

func (Tool) Name() string {
	return "Python"
}

func (t Tool) Executable() string {
	return "python3"
}

func (Tool) Filenames() []string {
	return nil
}

func (Tool) Extensions() []string {
	return []string{"py"}
}

func (Tool) ProcessFile(path string) ([]tool.Task, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read python file")
	}

	python := string(content)
	return lo.FilterMap(strings.Split(python, "\n"), func(line string, _ int) (tool.Task, bool) {
		captures := util.NamedCaptureGroups(definitionRegexp, line)
		if len(captures) == 0 {
			return nil, false
		}

		name := captures["name"]
		paramList := captures["paramList"]
		var params parameter.Parameters
		for _, param := range util.SplitCommaParamList(paramList) {
			fields := strings.SplitN(param, "=", 2)
			switch len(fields) {
			case 1:
				paramName := strings.TrimSpace(fields[0])
				params = append(params, parameter.NewRequired(paramName, parameter.TypeString))

			case 2:
				paramName, paramDefault := strings.TrimSpace(fields[0]), strings.TrimSpace(fields[1])
				params = append(params, parameter.NewOptional(paramName, parseType(paramDefault), paramDefault))

			default:
				logger.Log.Printf("invalid parameter in %q: %s", name, param)
			}
		}

		return Task{name: name, parameters: params, pythonCode: &python}, true
	}), nil
}

func parseType(param string) parameter.Type {
	if param == "False" || param == "True" {
		return parameter.TypeBool
	} else if lo.Every([]rune("1234567890_"), []rune(param)) {
		return parameter.TypeInt
	} else if lo.Every([]rune("1234567890_."), []rune(param)) {
		return parameter.TypeFloat
	} else if (strings.HasPrefix(param, `"`) && strings.HasSuffix(param, `"`)) || (strings.HasPrefix(param, `'`) && strings.HasSuffix(param, `'`)) {
		return parameter.TypeString
	} else {
		panic(fmt.Sprintf("invalid type %q", param))
	}
}
