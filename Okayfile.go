package main

import (
	"io"
	"os"
	"os/exec"

	"github.com/creack/pty"
)

func pry() {
	tty, err := pty.Start(exec.Command("sh"))
	if err != nil {
		panic(err)
	}

	defer tty.Close()

	go func() { _, _ = io.Copy(tty, os.Stdin) }()
	_, _ = io.Copy(os.Stdout, tty)
}
