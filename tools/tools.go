package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/broothie/ok/tool"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type Tools map[string]tool.Tool

func NewTools(tools []tool.Tool) Tools {
	return lo.Associate(tools, func(tool tool.Tool) (string, tool.Tool) { return tool.Name(), tool })
}

func FromRegistry() Tools {
	return NewTools(lo.Map(Registry(), func(newFunc tool.NewFunc, _ int) tool.Tool { return newFunc() }))
}

func (t Tools) CollectTasks() (Tasks, error) {
	tasks := make(map[string]Task)

	for _, tool := range t {
		paths := tool.Config().Filenames()
		for _, extension := range tool.Config().Extensions() {
			paths = append(paths, fmt.Sprintf("Okfile.%s", extension))

			subOkFiles, err := filepath.Glob(fmt.Sprintf("okfiles/*.%s", extension))
			if err != nil {
				return nil, fmt.Errorf("failed to glob for extension %q", extension)
			}

			paths = append(paths, subOkFiles...)
		}

		for _, filename := range paths {
			if _, err := os.Stat(filename); err != nil {
				if os.IsNotExist(err) {
					continue
				}

				return nil, errors.Wrap(err, "failed to stat file")
			}

			toolTasks, err := tool.ProcessFile(filename)
			if err != nil {
				return nil, err
			}

			prefix := ""
			basename := filepath.Base(filename)
			if filename != basename {
				prefix = fmt.Sprintf("%s.", strings.TrimSuffix(basename, filepath.Ext(basename)))
			}

			for _, task := range toolTasks {
				tasks[fmt.Sprintf("%s%s", prefix, task.Name())] = Task{Task: task, Tool: tool, Filename: filename}
			}
		}
	}

	return tasks, nil
}
