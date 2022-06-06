package install

import (
	"io"
)

type commandWriter struct {
	instr   *installer
	command *Command
}

func (w *commandWriter) Write(p []byte) (n int, err error) {
	w.command.Out = append(w.command.Out, p...)
	w.instr.triggerProgess()
	return len(p), nil
}

func newCommandWriter(instr *installer, command *Command) io.Writer {
	command.Out = []byte{}

	return &commandWriter{
		instr:   instr,
		command: command,
	}
}
