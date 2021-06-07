package task

type Task interface {
	Name() string
	Comment() string
	Filename() string
	Params() Parameters
	Invoke(args Args) (RunningTask, error)
}

type RunningTask interface {
	Wait() error
	Kill() error
}
