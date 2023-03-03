package make

import (
	"os"
	"regexp"
	"strings"

	"github.com/broothie/ok/task"
	"github.com/broothie/ok/tool"
	"github.com/broothie/ok/util"
	"github.com/pkg/errors"
)

var ruleRegexp = regexp.MustCompile(`^(?P<name>\w[a-zA-Z0-9]+):.*$`)

type Tool struct {
	config *tool.Config
}

func New() tool.Tool {
	return Tool{
		config: &tool.Config{
			"filenames":  "Makefile",
			"executable": "make",
		},
	}
}

func (Tool) Name() string {
	return "make"
}

func (t Tool) Config() *tool.Config {
	return t.config
}

func (t Tool) ProcessFile(path string) ([]task.Task, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read Makefile")
	}

	makeCode := string(content)
	var tasks []task.Task
	for _, line := range strings.Split(makeCode, "\n") {
		captures := util.NamedCaptureGroups(ruleRegexp, line)
		if len(captures) == 0 {
			continue
		}

		tasks = append(tasks, Task{
			Tool: t,
			name: captures["name"],
		})
	}

	return tasks, nil
}
