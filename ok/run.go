package ok

import (
	"fmt"

	"github.com/broothie/ok/cli"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

type Ok struct {
	Parser   *cli.Parser
	Options  cli.Options
	TaskList []Task
}

func New(args []string) (*Ok, error) {
	parser, err := cli.NewParser(args)
	if err != nil {
		return nil, err
	}

	return &Ok{Parser: parser}, nil
}

func Run(args []string) error {
	ok, err := New(args)
	if err != nil {
		return err
	}

	return ok.Run()
}

func (ok *Ok) Run() error {
	taskName, halt, err := ok.HandleOptions()
	if err != nil {
		return err
	} else if halt {
		return nil
	}

	mountErrors := ok.Mount()
	foundTask := funk.Find(ok.TaskList, func(task Task) bool { return task.Name() == taskName })
	if foundTask == nil {
		return fmt.Errorf("no task found with name '%s'", taskName)
	}

	task := foundTask.(Task)
	if err, errPresent := mountErrors[task.Tool.Name()]; errPresent {
		return errors.Wrapf(err, "failed to mount tool '%s'", task.Tool.Name())
	}

	taskArgs, err := ok.Parser.ParseArgs(task.Params())
	if err != nil {
		return err
	}

	if len(ok.Options.Watches) != 0 {
		return ok.runWatcher(task, taskArgs)
	}

	return task.Invoke(taskArgs).Wait()
}
