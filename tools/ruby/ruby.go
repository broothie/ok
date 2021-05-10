package ruby

import "os"

const (
	ToolName = "ruby"
	filename = "Okayfile.rb"
)

type Ruby struct{}

func (Ruby) Init() error {
	_, err := os.Create(filename)
	return err
}
