package ok

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bmatcuk/doublestar"
	"github.com/broothie/ok/logger"
	"github.com/broothie/ok/task"
	"github.com/pkg/errors"
	"github.com/radovskyb/watcher"
)

func (ok *Ok) runWatcher(t task.Task, args task.Args) error {
	watcher := watcher.New()
	watcher.SetMaxEvents(1)

	for _, watchPattern := range ok.Options.Watches {
		filenames, err := doublestar.Glob(watchPattern)
		if err != nil {
			return errors.Wrapf(err, "failed to glob %q", watchPattern)
		}

		for _, filename := range filenames {
			if err := watcher.Add(filename); err != nil {
				return errors.Wrapf(err, "failed to add file %q to watches", filename)
			}
		}
	}

	go func() {
		process, err := t.Invoke(args)
		if err != nil {
			logger.Ok.Printf("failed to start task process: %v", err)
		}

		for {
			select {
			case <-watcher.Event:
				if process != nil {
					if err := process.Kill(); err != nil {
						logger.Ok.Printf("failed to kill process: %v", err)
					}
				}

				process, err = t.Invoke(args)
				if err != nil {
					logger.Ok.Printf("failed to start task process: %v", err)
				}

			case err := <-watcher.Error:
				if process != nil {
					if err := process.Kill(); err != nil {
						logger.Ok.Printf("failed to kill process: %v", err)
					}
				}

				logger.Ok.Printf("watcher error: %v", err)

			case <-watcher.Closed:
				if process != nil {
					if err := process.Kill(); err != nil {
						logger.Ok.Printf("failed to kill process: %v", err)
					}
				}

				return
			}
		}
	}()

	kill := make(chan os.Signal)
	signal.Notify(kill, syscall.SIGINT)
	signal.Notify(kill, syscall.SIGTERM)
	go func() {
		<-kill
		watcher.Close()
	}()

	if err := watcher.Start(100 * time.Millisecond); err != nil {
		return errors.Wrap(err, "failed to start watcher")
	}

	return nil
}
