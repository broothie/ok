package ez

import (
	"bytes"
	"fmt"
	"os"
	"regexp"

	"github.com/broothie/ok/task"
	"github.com/broothie/ok/toolhelp"
	"github.com/pkg/errors"
)

type ParamParser func(paramString string) (task.Parameters, error)

type InvokeFunc func(task Task, args task.Args) *os.Process

type Tool struct {
	ToolName             string
	CommandName          string
	ToolFilename         string
	ToolInitContent      string
	TaskMatcher          *regexp.Regexp
	CommentPrefixMatcher *regexp.Regexp
	ParamParser          ParamParser
	Invoke               InvokeFunc
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
			return fmt.Errorf("file '%s' already exists", t.ToolFilename)
		}

		return errors.Wrapf(err, "failed to create file '%s'", t.ToolFilename)
	}

	if _, err := fmt.Fprint(file, t.ToolInitContent); err != nil {
		return errors.Wrapf(err, "could not write to file '%s'", t.ToolFilename)
	}

	return errors.Wrapf(file.Close(), "could not close file '%s'", t.ToolFilename)
}

func (t Tool) Check() error {
	return toolhelp.Check(t.CommandName)
}

func (t Tool) Mount() ([]task.Task, error) {
	fileBytes, err := os.ReadFile(t.ToolFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, toolhelp.ReadToolFileError{Filename: t.ToolFilename, Err: err}
	}

	if err := t.Check(); err != nil {
		return nil, toolhelp.CommandNotFoundError{CommandName: t.CommandName}
	}

	fileContents := string(fileBytes)
	rawTasks := toolhelp.Scan(bytes.NewReader(fileBytes), t.TaskMatcher, t.CommentPrefixMatcher)
	tasks := make([]task.Task, len(rawTasks))
	for i, rawTask := range rawTasks {
		taskName := rawTask.MatchData["taskName"]
		paramString := rawTask.MatchData["params"]

		params, err := t.ParamParser(paramString)
		if err != nil {
			toolhelp.Warn(t.ToolName, "")
			continue
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

	return tasks, nil
}

type Task struct {
	Tool         *Tool
	TaskName     string
	TaskComment  string
	TaskFilename string
	TaskParams   task.Parameters
	FileContents *string
}

func (t Task) Name() string {
	return t.TaskName
}

func (t Task) Comment() string {
	return t.TaskComment
}

func (t Task) Filename() string {
	return t.TaskFilename
}

func (t Task) ToolName() string {
	return t.Tool.ToolName
}

func (t Task) Params() task.Parameters {
	return t.TaskParams
}

func (t Task) Invoke(args task.Args) *os.Process {
	return t.Tool.Invoke(t, args)
}
