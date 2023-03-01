package tool

type Tool interface {
	Name() string
	Executable() string
	Filenames() []string
	Extensions() []string
	ProcessFile(path string) ([]Task, error)
}
