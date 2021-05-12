package make

import (
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"

	"github.com/broothie/ok/task"
	"github.com/broothie/ok/toolhelp"
)

var ruleMatcher = regexp.MustCompile(`(?m)^\s*(.*):`)

func (t Tool) Mount() ([]task.Task, error) {
	if _, err := os.Stat(filename); err != nil {
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

	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, toolhelp.ReadToolFileError{Filename: filename, Err: err}
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
