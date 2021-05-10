package make

import (
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"

	"github.com/broothie/now/task"
	"github.com/broothie/now/toolhelp"
)

var ruleMatcher = regexp.MustCompile(`(?m)^\s*(.*):`)

func (Make) Mount() ([]task.Task, error) {
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, toolhelp.ReadToolFileError{Filename: filename, Err: err}
	}

	if _, err := exec.LookPath(ToolName); err != nil {
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
			Base: task.NewBaseTask(taskName, filename, ToolName),
		}

		counter++
	}

	return tasks, nil
}
