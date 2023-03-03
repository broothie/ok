package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func (t Tools) CollectTasks() (Tasks, error) {
	tasks := make(Tasks)
	var tasksLock sync.Mutex

	group, ctx := errgroup.WithContext(context.Background())
	for _, tool := range t {
		tool := tool
		group.Go(func() error {
			paths := tool.Config().Filenames()
			for _, extension := range tool.Config().Extensions() {
				if ctx.Err() != nil {
					return nil
				}

				subOkFiles, err := filepath.Glob(fmt.Sprintf("okfiles/*.%s", extension))
				if err != nil {
					return fmt.Errorf("failed to glob for extension %q", extension)
				}

				paths = append(paths, subOkFiles...)
				paths = append(paths, fmt.Sprintf("Okfile.%s", extension))
			}

			group.Go(func() error {
				for _, filename := range paths {
					if ctx.Err() != nil {
						return nil
					}

					if _, err := os.Stat(filename); err != nil {
						if os.IsNotExist(err) {
							return nil
						}

						return errors.Wrap(err, "failed to stat file")
					}

					toolTasks, err := tool.ProcessFile(filename)
					if err != nil {
						return err
					}

					prefix := ""
					basename := filepath.Base(filename)
					if filename != basename {
						prefix = fmt.Sprintf("%s.", strings.TrimSuffix(basename, filepath.Ext(basename)))
					}

					for _, task := range toolTasks {
						name := fmt.Sprintf("%s%s", prefix, task.Name())

						tasksLock.Lock()
						tasks[name] = Task{
							Task:     task,
							Tool:     tool,
							Filename: filename,
						}
						tasksLock.Unlock()
					}
				}

				return nil
			})

			return nil
		})
	}

	if err := group.Wait(); err != nil {
		return nil, err
	}

	return tasks, nil
}
