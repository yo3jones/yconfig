package install

import (
	"io"
	"os/exec"
)

func ExecBashCommand(command string, out io.Writer) error {
	var (
		stdout io.Reader
		stderr io.Reader
		err    error
	)

	cmd := exec.Command("bash", "-c", command)

	if stdout, err = cmd.StdoutPipe(); err != nil {
		return err
	}
	if stderr, err = cmd.StderrPipe(); err != nil {
		return err
	}

	if err = cmd.Start(); err != nil {
		return err
	}

	if _, err = io.Copy(out, stdout); err != nil {
		return err
	}

	if _, err = io.Copy(out, stderr); err != nil {
		return err
	}

	if err = cmd.Wait(); err != nil {
		return err
	}

	return nil
}
