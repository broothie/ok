package ok

import "context"

type Tool struct {
	Name        string
	CommandName string
	FileGlobs   []string
	ParseFile   func(ctx context.Context, filePath string) ([]Task, error)
}

type Task struct {
	Name string
	Run  func(ctx context.Context, args []string) error
}
