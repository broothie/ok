package runner

import (
	"fmt"
	"os"

	"github.com/broothie/okay/tool"

	"github.com/broothie/okay/arg"
	"github.com/broothie/okay/task"
)

func Run(args []string) {
	New(args).Run()
}

type Runner struct {
	arg.Parser
	Tasks map[string]task.Task
	Task  task.Task
}

func New(argList []string) Runner {
	return Runner{
		Parser: arg.NewParser(argList),
		Tasks:  make(map[string]task.Task),
	}
}

func (r Runner) Run() {
	if err := r.ParseOptions(); err != nil {
		log("error parsing args: %v", err)
		os.Exit(1)
		return
	}

	if !r.Options.WillRunTask() {
		if r.Options.Help {
			arg.PrintHelp()
		} else if r.Options.ListTools {
			// TODO: List availability states
			for toolName := range tool.Registry {
				fmt.Println(toolName)
			}
		} else if r.Options.Init != "" {
			toolName := r.Options.Init
			tool, ok := tool.Registry[r.Options.Init]
			if !ok {
				log("no tool with name '%s'", toolName)
				os.Exit(1)
				return
			}

			if err := tool.Init(); err != nil {
				log("error initializing tool '%s'", toolName)
			}
		} else {
			r.mount()
			if err := r.list(); err != nil {
				log("error listing tasks: %v", err)
				os.Exit(1)
				return
			}
		}

		os.Exit(0)
		return
	}

	r.mount()
	taskName := r.Options.TaskName
	if task, taskExists := r.Tasks[taskName]; !taskExists {
		log("no task with name '%s'", taskName)
		os.Exit(1)
		return
	} else {
		r.Task = task
	}

	if err := r.ParseTaskArgs(r.Task.Params()); err != nil {
		log(err.Error())
		os.Exit(1)
		return
	}

	if len(r.Options.Watches) > 0 {
		r.handleWatches()
	} else {
		r.Task.Invoke(r.Parser.Args).Wait()
	}
}
