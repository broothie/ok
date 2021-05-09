package driver

import (
	"fmt"
	"os"

	"github.com/broothie/now/arg"
	"github.com/broothie/now/task"
)

func Run(args []string) {
	New(args).Run()
}

type Driver struct {
	arg.Parser
	Tasks map[string]task.Task
	Task  task.Task
}

func New(argList []string) Driver {
	return Driver{
		Parser: arg.NewParser(argList),
		Tasks:  make(map[string]task.Task),
	}
}

func (d Driver) Run() {
	if err := d.ParseNowArgs(); err != nil {
		log("error parsing args: %v", err)
		os.Exit(1)
		return
	}

	if !d.NowArgs.WillRunTask() {
		if d.NowArgs.Help {
			fmt.Println("help")
		} else if d.NowArgs.Explain {
			fmt.Println("explain")
		} else {
			d.mount()
			if err := d.list(); err != nil {
				log("error listing tasks: %v", err)
				os.Exit(1)
				return
			}
		}

		os.Exit(0)
		return
	}

	d.mount()
	taskName := d.NowArgs.TaskName
	if task, taskExists := d.Tasks[taskName]; !taskExists {
		log("no task with name '%s'", taskName)
		os.Exit(1)
		return
	} else {
		d.Task = task
	}

	if err := d.ParseTaskArgs(d.Task.Params()); err != nil {
		log(err.Error())
		os.Exit(1)
		return
	}

	if len(d.NowArgs.Watches) > 0 {
		d.handleWatches()
	} else {
		d.Task.Invoke(d.Parser.TaskArgs).Wait()
	}
}
