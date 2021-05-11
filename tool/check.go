package tool

import "os/exec"

func Check(command string) error {
	if _, err := exec.LookPath(command); err != nil {
		return CommandNotFoundError{CommandName: command}
	}

	return nil
}
