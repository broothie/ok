package dockercompose

import (
	"os"
	"os/exec"

	"github.com/broothie/ok/task"
	"github.com/broothie/ok/tool"
	"gopkg.in/yaml.v3"
)

func (t Tool) Mount() ([]task.Task, error) {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, tool.ReadToolFileError{Filename: filename, Err: err}
	}

	defer file.Close()

	if err := t.Check(); err != nil {
		if err == exec.ErrNotFound {
			return nil, tool.CommandNotFoundError{CommandName: ToolName}
		}

		return nil, err
	}

	var dockerComposeYML map[string]interface{}
	if err := yaml.NewDecoder(file).Decode(&dockerComposeYML); err != nil {
		return nil, tool.ReadToolFileError{Filename: filename, Err: err}
	}

	untypedServices, servicesPresent := dockerComposeYML["services"]
	if !servicesPresent {
		return nil, nil
	}

	services := untypedServices.(map[string]interface{})
	tasks := make([]task.Task, len(services))
	counter := 0
	for serviceName := range services {
		tasks[counter] = Task{Base: task.NewBase(serviceName, filename, ToolName)}
		counter++
	}

	return tasks, nil
}
