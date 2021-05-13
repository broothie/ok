package rake

import (
	"os"
	"os/exec"
	"regexp"

	"github.com/broothie/ok/stringhelp"
	"github.com/broothie/ok/task"
	"github.com/broothie/ok/toolhelp"
	"github.com/pkg/errors"
)

var taskMatcher = regexp.MustCompile(`(?m)^rake (?P<taskName>\w+)(?:\[(?P<params>.*?)])?\s+# (?P<description>.*)?$`)

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

	taskList, err := exec.Command(ToolName, "-AT").Output()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get rake tasks")
	}

	results := stringhelp.NamedRegexpResults(string(taskList), taskMatcher)
	tasks := make([]task.Task, len(results))
	for i, result := range results {
		taskName := result["taskName"]
		paramsString := result["params"]
		description := result["description"]

		paramStrings := stringhelp.SplitOnCommas(paramsString)
		paramList := make(task.ParamList, len(paramStrings))
		for i, paramString := range paramStrings {
			paramList[i] = task.Parameter{Name: paramString, Type: task.String}
		}

		tasks[i] = Task{
			Base:    task.NewBase(taskName, filename, ToolName),
			params:  paramList.ToParameters(false),
			comment: description,
		}
	}

	return tasks, nil
}
