package install

import (
	"io"
	"os"
	"os/exec"
)

func ExecBashCommand(command string, out io.Writer) error {
	var err error

	cmd := exec.Command("bash", "-c", command)

	cmd.Env = os.Environ()

	cmd.Stdout = out
	cmd.Stderr = out

	if err = cmd.Start(); err != nil {
		return err
	}

	if err = cmd.Wait(); err != nil {
		return err
	}

	return nil
}
