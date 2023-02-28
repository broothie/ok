package ruby

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

var definitionRegexp = regexp.MustCompile(`^def (?P<name>\w[a-zA-Z0-9_]*)(?:\(?(?P<paramList>[^)]*)\)?)?$`)

type Tool struct{}

func (Tool) Name() string {
	return "ruby"
}

func (Tool) CommandName() string {
	return "ruby"
}

func (Tool) Filenames() []string {
	return nil
}

func (Tool) Extensions() []string {
	return []string{"rb"}
}

func (Tool) ProcessFile(path string) ([]tool.Task, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read ruby file")
	}

	ruby := string(content)

	return lo.FilterMap(strings.Split(ruby, "\n"), func(line string, _ int) (tool.Task, bool) {
		captures := util.NamedCaptureGroups(definitionRegexp, line)
		if len(captures) == 0 {
			return nil, false
		}

		name := captures["name"]
		paramList := captures["paramList"]
		var params parameter.Parameters
		for _, arg := range util.SplitCommaParamList(paramList) {
			fields := strings.Fields(arg)
			switch len(fields) {
			case 1:
				name := fields[0]
				params = append(params, parameter.Parameter{
					Name: name,
					Type: parameter.TypeString,
				})

			case 2:
				name, dflt := fields[0], fields[1]
				params = append(params, parameter.Parameter{
					Name:    strings.TrimSuffix(name, ":"),
					Type:    parseType(dflt),
					Default: &dflt,
				})

			default:
				logger.Log.Println("invalid parameter:", arg)
			}
		}

		return Task{name: name, parameters: params}, true
	}), nil
}

func parseType(param string) parameter.Type {
	if param == "false" || param == "true" {
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
