package golang

import (
	"fmt"
	"os"
)

const (
	ToolName = "go"
	filename = "Okayfile.go"
)

type Golang struct{}

func (Golang) Init() error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = fmt.Fprint(file, "package main\n")
	return err
}
