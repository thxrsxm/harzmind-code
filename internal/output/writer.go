package output

import (
	"fmt"
	"io"
)

// outWriter wraps an io.Writer to provide convenient print methods.
type outWriter struct {
	writer io.Writer
}

// newOutWriter creates a new outWriter instance with the given io.Writer.
func newOutWriter(writer io.Writer) *outWriter {
	return &outWriter{writer: writer}
}

// print writes the provided values to the underlying writer without formatting (no newline).
func (w *outWriter) print(a ...any) {
	if w.writer == nil {
		return
	}
	fmt.Fprint(w.writer, a...)
}

// println writes the provided values to the underlying writer with a newline.
func (w *outWriter) println(a ...any) {
	if w.writer == nil {
		return
	}
	fmt.Fprintln(w.writer, a...)
}

// printf writes a formatted string to the underlying writer using fmt.Printf format.
func (w *outWriter) printf(format string, a ...any) {
	if w.writer == nil {
		return
	}
	fmt.Fprintf(w.writer, format, a...)
}
