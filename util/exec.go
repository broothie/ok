package util

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/pkg/errors"
)

type ExecProcess struct {
	Cmd           *exec.Cmd
	killChan      chan struct{}
	killErrorChan chan error
}

func (p *ExecProcess) Kill() error {
	close(p.killChan)
	return <-p.killErrorChan
}

func (p *ExecProcess) Wait() error {
	return p.Cmd.Wait()
}

func Exec(name string, arg ...string) (*ExecProcess, error) {
	cmd := exec.Command(name, arg...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, errors.Wrap(err, "failed to start command")
	}

	process := &ExecProcess{
		Cmd:           cmd,
		killChan:      make(chan struct{}),
		killErrorChan: make(chan error),
	}

	go func() {
		<-process.killChan
		process.killErrorChan <- syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}()

	return process, nil
}
