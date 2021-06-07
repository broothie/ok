package util

import (
	"os"
	"os/exec"
	"sync"
	"syscall"

	"github.com/pkg/errors"
)

type ExecProcess struct {
	*exec.Cmd
	killed    bool
	waiter    *sync.WaitGroup
	waitError error
}

func (p ExecProcess) Kill() error {
	p.killed = true
	return syscall.Kill(-p.Process.Pid, syscall.SIGKILL)
}

func (p ExecProcess) Wait() error {
	p.waiter.Wait()
	return p.waitError
}

func Exec(name string, arg ...string) (ExecProcess, error) {
	cmd := exec.Command(name, arg...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return ExecProcess{}, errors.Wrap(err, "")
	}

	process := ExecProcess{Cmd: cmd, waiter: new(sync.WaitGroup)}
	process.waiter.Add(1)
	go func() {
		defer process.waiter.Done()

		if process.killed || (process.ProcessState != nil && process.ProcessState.Exited()) {
			return
		}

		process.waitError = cmd.Wait()
	}()

	return process, nil
}
