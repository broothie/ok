package ok

import (
	"fmt"

	"github.com/broothie/ok/cli"
	"github.com/broothie/ok/config"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

type Ok struct {
	Version   string
	Parser    *cli.Parser
	Options   cli.Options
	MapConfig map[string]interface{}
	TaskList  []Task
}

func New(version string, args []string) (*Ok, error) {
	var cfg config.Config
	config.ReadConfigAndEnv(&cfg)
	parser, err := cli.NewParser(args, cfg)
	if err != nil {
		return nil, err
	}

	mapConfig := make(map[string]interface{})
	config.ReadInConfig(&mapConfig)
	return &Ok{Parser: parser, MapConfig: mapConfig}, nil
}

func Run(version string, args []string) error {
	ok, err := New(version, args)
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
		return fmt.Errorf("no task found with name %q", taskName)
	}

	task := foundTask.(Task)
	if err, errPresent := mountErrors[task.Tool.Name()]; errPresent {
		return errors.Wrapf(err, "failed to mount tool %q", task.Tool.Name())
	}

	if ok.Options.Debug {
		fmt.Printf("%+v", task)
	}

	taskArgs, err := ok.Parser.ParseArgs(task.Params())
	if err != nil {
		return err
	}

	if len(ok.Options.Watches) != 0 {
		return ok.runWatcher(task, taskArgs)
	}

	process, err := task.Invoke(taskArgs)
	if err != nil {
		return errors.Wrap(err, "failed to start task process")
	}

	return process.Wait()
}
