package yarn

import (
	"encoding/json"
	"os"
	"os/exec"

	"github.com/broothie/okay/task"
	"github.com/broothie/okay/toolhelp"
)

func (Yarn) Mount() ([]task.Task, error) {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, toolhelp.ReadToolFileError{Filename: filename, Err: err}
	}

	defer file.Close()

	if _, err := exec.LookPath("yarn"); err != nil {
		if err == exec.ErrNotFound {
			return nil, toolhelp.CommandNotFoundError{CommandName: ToolName}
		}

		return nil, err
	}

	var packageJSON map[string]interface{}
	if err := json.NewDecoder(file).Decode(&packageJSON); err != nil {
		return nil, toolhelp.ReadToolFileError{Filename: filename, Err: err}
	}

	untypedScripts, scriptsPresent := packageJSON["scripts"]
	if !scriptsPresent {
		return nil, nil
	}

	scripts := untypedScripts.(map[string]interface{})
	tasks := make([]task.Task, len(scripts))
	counter := 0
	for scriptName := range scripts {
		tasks[counter] = Task{
			Base: task.NewBaseTask(scriptName, filename, ToolName),
		}

		counter++
	}

	return tasks, nil
}
