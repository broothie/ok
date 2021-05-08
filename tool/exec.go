package tool

import (
	"os"
	"os/exec"
)

func Exec(name string, arg ...string) *exec.Cmd {
	command := exec.Command(name, arg...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Start()

	return command
}
