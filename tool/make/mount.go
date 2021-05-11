package make

import (
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"

	"github.com/broothie/okay/task"
	"github.com/broothie/okay/tool"
)

var ruleMatcher = regexp.MustCompile(`(?m)^\s*(.*):`)

func (m Make) Mount() ([]task.Task, error) {
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, tool.ReadToolFileError{Filename: filename, Err: err}
	}

	if err := m.Check(); err != nil {
		if err == exec.ErrNotFound {
			return nil, tool.CommandNotFoundError{CommandName: ToolName}
		}

		return nil, err
	}

	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, tool.ReadToolFileError{Filename: filename, Err: err}
	}

	matches := ruleMatcher.FindAllStringSubmatch(string(fileBytes), -1)
	tasks := make([]task.Task, len(matches))
	counter := 0
	for _, match := range matches {
		taskName := match[1]
		tasks[counter] = Task{
			Base: task.NewBase(taskName, filename, ToolName),
		}

		counter++
	}

	return tasks, nil
}
