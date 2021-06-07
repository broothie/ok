package ez

import "github.com/broothie/ok/task"

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

func (t Task) Params() task.Parameters {
	return t.TaskParams
}

func (t Task) Invoke(args task.Args) (task.RunningTask, error) {
	return t.Tool.Invoke(t, args)
}
