package setup

import "io"

type setupWriter struct {
	buffer     *[]byte
	onProgress func()
}

func (w *setupWriter) Write(p []byte) (n int, err error) {
	*w.buffer = append(*w.buffer, p...)
	w.onProgress()
	return len(p), nil
}

func NewWriter(out *[]byte, onProgress func()) (writer io.Writer) {
	return &setupWriter{
		buffer:     out,
		onProgress: onProgress,
	}
}
