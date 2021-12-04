package util

import (
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/broothie/ok/logger"
	"github.com/pkg/errors"
)

type ExecProcess struct {
	*exec.Cmd
	killChan      chan struct{}
	killErrorChan chan error
}

func (p *ExecProcess) Kill() error {
	close(p.killChan)
	return <-p.killErrorChan
}

func Exec(name string, arg ...string) (*ExecProcess, error) {
	cmd := exec.Command(name, arg...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get stdin pipe")
	}

	go func() {
		if _, err := io.Copy(stdin, os.Stdin); err != nil {
			logger.Ok.Printf("failed to pipe stdin: %v", err)
		}
	}()

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
		process.killErrorChan <- syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)
	}()

	return process, nil
}
