package make

import (
	"os"
	"os/exec"
	"regexp"

	"github.com/broothie/ok/stringhelp"
	"github.com/broothie/ok/task"
	"github.com/broothie/ok/toolhelp"
)

var (
	ruleMatcher = regexp.MustCompile(`(?m)^\s*(?P<taskName>.*?):`)
)

func (t Tool) Mount() ([]task.Task, error) {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, toolhelp.ReadToolFileError{Filename: filename, Err: err}
	}

	if err := t.Check(); err != nil {
		if err == exec.ErrNotFound {
			return nil, toolhelp.CommandNotFoundError{CommandName: ToolName}
		}

		return nil, err
	}

	rawTasks := toolhelp.Scan(file, ruleMatcher, stringhelp.OctothorpePrefixMatcher)
	tasks := make([]task.Task, len(rawTasks))
	for i, rawTask := range rawTasks {
		taskName := rawTask.MatchData["taskName"]

		tasks[i] = Task{
			Base:    task.NewBase(taskName, filename, ToolName),
			comment: rawTask.Comment,
		}
	}

	return tasks, nil
}
