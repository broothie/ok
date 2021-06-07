package util

import (
	"os"
	"os/exec"
	"syscall"
)

type ExecProcess struct {
	*exec.Cmd
}

func (p ExecProcess) Kill() error {
	return syscall.Kill(-p.Process.Pid, syscall.SIGKILL)
}

func Exec(name string, arg ...string) (ExecProcess, error) {
	cmd := exec.Command(name, arg...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return ExecProcess{Cmd: cmd}, cmd.Start()
}
