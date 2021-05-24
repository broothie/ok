package toolhelp

import (
	"os"
	"os/exec"
	"syscall"
)

type ExecProcess struct {
	*exec.Cmd
}

func (c ExecProcess) Kill() error {
	return syscall.Kill(-c.Process.Pid, syscall.SIGKILL)
}

func Exec(name string, arg ...string) ExecProcess {
	cmd := exec.Command(name, arg...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()

	return ExecProcess{Cmd: cmd}
}
