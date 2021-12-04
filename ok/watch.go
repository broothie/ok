package ok

import (
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/bmatcuk/doublestar"
	"github.com/broothie/ok/logger"
	"github.com/broothie/ok/task"
	"github.com/pkg/errors"
	"github.com/radovskyb/watcher"
)

func (ok *Ok) runWatcher(t task.Task, args task.Args) error {
	var wg sync.WaitGroup
	defer wg.Wait()
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

	wg.Add(1)
	go func() {
		defer wg.Done()
		var process task.RunningTask

		for {
			select {
			case <-watcher.Event:
				if process != nil {
					process.Kill()
				}

				var err error
				process, err = t.Invoke(args)
				if err != nil {
					logger.Ok.Printf("failed to start task process: %v", err)
				}

			case err := <-watcher.Error:
				if process != nil {
					process.Kill()
				}

				logger.Ok.Println(err)

			case <-watcher.Closed:
				if process != nil {
					process.Kill()
				}

				return
			}
		}
	}()

	kill := make(chan os.Signal)
	signal.Notify(kill, os.Interrupt)
	signal.Notify(kill, os.Kill)
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-kill
		watcher.Close()
	}()

	if err := watcher.Start(100 * time.Millisecond); err != nil {
		return errors.Wrap(err, "failed to start watcher")
	}

	return nil
}
