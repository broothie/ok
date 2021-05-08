package now

import (
	"os"
	"path/filepath"
	"time"

	"github.com/radovskyb/watcher"
)

func (n Now) handleWatches() {
	watcher := watcher.New()
	watcher.SetMaxEvents(1)
	for _, watchPattern := range n.NowArgs.Watches {
		filePaths, err := filepath.Glob(watchPattern)
		if err != nil {
			log("error globbing glob '%s': %v", watchPattern, err)
			os.Exit(1)
			return
		}

		for _, filePath := range filePaths {
			if err := watcher.Add(filePath); err != nil {
				log("error adding watch '%s': %v", filePath, err)
				os.Exit(1)
				return
			}
		}
	}

	go func() {
		var process *os.Process
		defer func() {
			if process != nil {
				process.Kill()
			}
		}()

		for {
			select {
			case <-watcher.Event:
				if process != nil {
					process.Kill()
					process = nil
				}

				process = n.Task.Invoke(n.Parser.TaskArgs)

			case <-watcher.Closed:
				break
			}
		}
	}()

	if err := watcher.Start(100 * time.Millisecond); err != nil {
		log("error starting watch process: %v\n", err)
		os.Exit(1)
		return
	}
}
