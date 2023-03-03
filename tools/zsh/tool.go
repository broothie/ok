package zsh

import (
	"os"
	"regexp"
	"strings"

	"github.com/broothie/ok/task"
	"github.com/broothie/ok/tool"
	"github.com/broothie/ok/util"
	"github.com/pkg/errors"
)

const commentPrefix = "#"

var definitionRegexp = regexp.MustCompile(`^function (?P<name>\w[a-zA-Z0-9_]*)\(\) \{$`)

type Tool struct {
	config *tool.Config
}

func New() tool.Tool {
	return Tool{
		config: &tool.Config{
			"extensions": "zsh",
			"executable": "zsh",
		},
	}
}

func (Tool) Name() string {
	return "zsh"
}

func (t Tool) Config() *tool.Config {
	return t.config
}

func (Tool) Init() error {
	return util.InitFile("Okfile.zsh", nil)
}

func (t Tool) ProcessFile(path string) ([]task.Task, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read zsh file")
	}

	zsh := string(content)
	lines := strings.Split(zsh, "\n")
	var tasks []task.Task
	for i, line := range lines {
		captures := util.NamedCaptureGroups(definitionRegexp, line)
		if len(captures) == 0 {
			continue
		}

		description := ""
		if i != 0 {
			description = util.ExtractComment(lines[i-1], commentPrefix)
		}

		tasks = append(tasks, Task{
			Tool:        t,
			name:        captures["name"],
			description: description,
			zshCode:     &zsh,
		})
	}

	return tasks, nil
}
