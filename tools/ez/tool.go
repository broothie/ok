package ez

import (
	"bytes"
	"fmt"
	"os"
	"regexp"

	"github.com/broothie/ok/task"
	"github.com/broothie/ok/util"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
)

type ParamParser func(paramString string) (task.Parameters, error)

type InvokeFunc func(task Task, args task.Args) (task.RunningTask, error)

type Tool struct {
	ToolName             string
	CommandName          string
	ToolFilename         string
	ToolInitContent      string
	TaskMatcher          *regexp.Regexp
	CommentPrefixMatcher *regexp.Regexp
	ParamParser          ParamParser
	Invoke               InvokeFunc
	ToolConfig           interface{}
}

func (t Tool) Name() string {
	return t.ToolName
}

func (t Tool) Filename() string {
	return t.ToolFilename
}

func (t Tool) Init() error {
	file, err := os.OpenFile(t.ToolFilename, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("file %q already exists", t.ToolFilename)
		}

		return errors.Wrapf(err, "failed to create file %q", t.ToolFilename)
	}

	if _, err := fmt.Fprint(file, t.ToolInitContent); err != nil {
		return errors.Wrapf(err, "could not write to file %q", t.ToolFilename)
	}

	return errors.Wrapf(file.Close(), "could not close file %q", t.ToolFilename)
}

func (t Tool) Check() error {
	return util.Check(t.CommandName)
}

func (t Tool) Configure(decoder *toml.Decoder) error {
	if t.ToolConfig == nil {
		return nil
	}

	return decoder.Decode(t.ToolConfig)
}

func (t Tool) Mount() ([]task.Task, error) {
	fileBytes, err := os.ReadFile(t.ToolFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, util.ReadToolFileError{Filename: t.ToolFilename, Err: err}
	}

	if err := t.Check(); err != nil {
		return nil, err
	}

	errs := new(util.ErrorGroup)
	fileContents := string(fileBytes)
	rawTasks := util.Scan(bytes.NewReader(fileBytes), t.TaskMatcher, t.CommentPrefixMatcher)
	tasks := make([]task.Task, len(rawTasks))
	for i, rawTask := range rawTasks {
		taskName := rawTask.MatchData["taskName"]
		paramString := rawTask.MatchData["params"]

		var params task.Parameters
		if t.ParamParser != nil {
			var err error
			if params, err = t.ParamParser(paramString); err != nil {
				errs.Add(errors.Wrapf(err, "failed to parse params for task %q from tool %q", taskName, t.Name()))
				continue
			}
		}

		tasks[i] = Task{
			Tool:         &t,
			TaskName:     taskName,
			TaskComment:  rawTask.Comment,
			TaskFilename: t.ToolFilename,
			TaskParams:   params,
			FileContents: &fileContents,
		}
	}

	return tasks, errs.NilIfEmpty()
}
