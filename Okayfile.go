package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/broothie/okay/okay"
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

func glob(glob string) {
	fmt.Println(filepath.Glob(glob))
}

func help() {
	p := okay.NewParser([]string{"-h"})
	p.WriteHelp(os.Stdout)
}

func oneChar() {
	s := "hello"
	fmt.Println(s[:1])
}
