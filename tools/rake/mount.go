package rake

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/broothie/ok/task"
	"github.com/broothie/ok/util"
)

var taskMatcher = regexp.MustCompile(`(?m)^rake (?P<taskName>\w+)(?:\[(?P<params>.*?)])?\s+# (?P<description>.*)?$`)

func (t *Tool) Mount() ([]task.Task, error) {
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, util.ReadToolFileError{Filename: filename, Err: err}
	}

	if err := t.Check(); err != nil {
		if err == exec.ErrNotFound {
			return nil, util.CommandNotFoundError{CommandName: ToolName}
		}

		return nil, err
	}

	command := "rake"
	args := []string{"-AT"}
	if t.ToolConfig.Bundle != nil && *t.ToolConfig.Bundle {
		command = "bundle"
		args = []string{"exec", "rake", "-AT"}
	}

	output, err := exec.Command(command, args...).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get rake tasks: %v: %s", err, output)
	}

	results := util.NamedRegexpResults(string(output), taskMatcher)
	tasks := make([]task.Task, len(results))
	for i, result := range results {
		taskName := result["taskName"]
		paramsString := result["params"]
		description := result["description"]

		paramStrings := util.SplitOnCommas(paramsString)
		paramList := make(task.ParamList, len(paramStrings))
		for i, paramString := range paramStrings {
			paramList[i] = task.Parameter{Name: paramString, Type: task.String}
		}

		tasks[i] = Task{
			Base:    task.NewBase(taskName, filename, ToolName),
			params:  paramList.ToParameters(false),
			comment: description,
			tool:    t,
		}
	}

	return tasks, nil
}
