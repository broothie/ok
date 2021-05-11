package golang

import (
	"fmt"
	"os"

	"github.com/broothie/okay/tool"
)

const (
	ToolName = "go"
	filename = "Okayfile.go"
)

type Tool struct{}

func (Tool) Name() string {
	return ToolName
}

func (Tool) Init() error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = fmt.Fprint(file, "package main\n")
	return err
}

func (Tool) Check() error {
	return tool.Check(ToolName)
}
