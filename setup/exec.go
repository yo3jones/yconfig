package setup

import (
	"io"
	"os"
	"os/exec"
)

func Exec(cmd string, args []string, writer io.Writer) (err error) {
	command := exec.Command(cmd, args...)

	command.Env = os.Environ()

	command.Stdout = writer
	command.Stderr = writer

	if err = command.Start(); err != nil {
		return err
	}

	if err = command.Wait(); err != nil {
		return err
	}

	return nil
}
