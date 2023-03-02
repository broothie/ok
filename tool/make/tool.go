package make

import (
	"os"
	"regexp"
	"strings"

	"github.com/broothie/ok/task"
	"github.com/broothie/ok/util"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

var ruleRegexp = regexp.MustCompile(`^(?P<name>\w[a-zA-Z0-9]+):.*$`)

type Tool struct{}

func (Tool) Name() string {
	return "Make"
}

func (Tool) Executable() string {
	return "make"
}

func (Tool) Filenames() []string {
	return []string{"Makefile"}
}

func (Tool) Extensions() []string {
	return nil
}

func (Tool) ProcessFile(path string) ([]task.Task, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read Makefile")
	}

	makeCode := string(content)
	return lo.FilterMap(strings.Split(makeCode, "\n"), func(line string, _ int) (task.Task, bool) {
		captures := util.NamedCaptureGroups(ruleRegexp, line)
		if len(captures) == 0 {
			return nil, false
		}

		name := captures["name"]
		return Task{name: name}, true
	}), nil
}
