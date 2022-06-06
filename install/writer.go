package install

import (
	"io"
)

type commandWriter struct {
	command *Command
}

func (w *commandWriter) Write(p []byte) (n int, err error) {
	w.command.Out = append(w.command.Out, p...)
	return len(p), nil
}

func newCommandWriter(command *Command) io.Writer {
	command.Out = []byte{}
	return &commandWriter{command}
}
