package now

import (
	"fmt"
	"os"

	"github.com/broothie/now/arg"
	"github.com/broothie/now/task"
)

func Run(args []string) {
	New(args).Run()
}

type Now struct {
	arg.Parser
	Tasks map[string]task.Task
	Task  task.Task
}

func New(argList []string) Now {
	return Now{
		Parser: arg.NewParser(argList),
		Tasks:  make(map[string]task.Task),
	}
}

func (n Now) Run() {
	n.ParseNowArgs()
	if !n.NowArgs.WillRunTask() {
		if n.NowArgs.Help {
			fmt.Println("help")
		} else if n.NowArgs.Explain {
			fmt.Println("explain")
		} else {
			n.mount()
			if err := n.list(); err != nil {
				log("error listing tasks: %v", err)
				os.Exit(1)
				return
			}
		}

		os.Exit(0)
		return
	}

	n.mount()
	taskName := n.NowArgs.TaskName
	if task, taskExists := n.Tasks[taskName]; !taskExists {
		log("no task with name '%s'", taskName)
		os.Exit(1)
		return
	} else {
		n.Task = task
	}

	if err := n.ParseTaskArgs(n.Task.Params()); err != nil {
		log(err.Error())
		os.Exit(1)
		return
	}

	if len(n.NowArgs.Watches) > 0 {
		n.handleWatches()
	} else {
		n.Task.Invoke(n.Parser.TaskArgs).Wait()
	}
}
