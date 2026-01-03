package output

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/thxrsxm/rnbw"
)

// OutputWriter represents a writer for output.
type OutputWriter struct {
	writer io.Writer
}

// NewOutputWriter creates a new OutputWriter instance.
func NewOutputWriter(writer io.Writer) *OutputWriter {
	return &OutputWriter{writer: writer}
}

// Print writes the given arguments to the output writer.
func (w *OutputWriter) Print(a ...any) {
	if w.writer == nil {
		return
	}
	fmt.Fprint(w.writer, a...)
}

// Println writes the given arguments to the output writer followed by a newline.
func (w *OutputWriter) Println(a ...any) {
	if w.writer == nil {
		return
	}
	fmt.Fprintln(w.writer, a...)
}

// Printf writes a formatted string to the output writer.
func (w *OutputWriter) Printf(format string, a ...any) {
	if w.writer == nil {
		return
	}
	fmt.Fprintf(w.writer, format, a...)
}

// Output represents a collection of output writers.
type Output struct {
	writer []OutputWriter
	Stdout *OutputWriter
	File   *OutputWriter
}

// NewOutput creates a new Output instance.
func NewOutput(outPath string, writeToFile bool) (*Output, error) {
	o := &Output{writer: []OutputWriter{}}
	// Add stdout writer to output
	o.writer = append(o.writer, *NewOutputWriter(os.Stdout))
	o.Stdout = &o.writer[0]
	// Prepare to write output to file
	if writeToFile {
		// Check out directory not exists
		if _, err := os.Stat(outPath); os.IsNotExist(err) {
			// Create out directory
			err := os.Mkdir(outPath, 0755)
			if err != nil {
				return nil, err
			}
		}
		// Set up out file with timestamp-based name
		outFilePath := fmt.Sprintf("%s/hzmind_%s.md", outPath, time.Now().Format("2006-01-02_15-04-05"))
		outFile, err := os.OpenFile(outFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		o.writer = append(o.writer, *NewOutputWriter(outFile))
		o.File = &o.writer[1]
	}
	return o, nil
}

// CloseOutput closes all output writers.
func (o *Output) CloseOutput() {
	for _, v := range o.writer {
		if closer, ok := v.writer.(io.Closer); ok {
			closer.Close()
		}
	}
}

// Print writes the given arguments to all output writers.
func (o *Output) Print(a ...any) {
	for _, v := range o.writer {
		v.Print(a...)
	}
}

// Println writes the given arguments to all output writers followed by a newline.
func (o *Output) Println(a ...any) {
	for _, v := range o.writer {
		v.Println(a...)
	}
}

// Printf writes a formatted string to all output writers.
func (o *Output) Printf(format string, a ...any) {
	for _, v := range o.writer {
		v.Printf(format, a...)
	}
}

// PrintWarning prints a warning message.
func (o *Output) PrintWarning(a ...any) {
	rnbw.ForgroundColor(rnbw.Yellow)
	o.Print("[WARNING] ")
	o.Print(a...)
	rnbw.ResetColor()
}

// PrintlnWarning prints a warning message followed by a newline.
func (o *Output) PrintlnWarning(a ...any) {
	rnbw.ForgroundColor(rnbw.Yellow)
	o.Print("[WARNING] ")
	o.Println(a...)
	rnbw.ResetColor()
}

// PrintfWarning prints a warning message with a formatted string.
func (o *Output) PrintfWarning(format string, a ...any) {
	rnbw.ForgroundColor(rnbw.Yellow)
	o.Print("[WARNING] ")
	o.Printf(format, a...)
	rnbw.ResetColor()
}

// PrintError prints an error message.
func (o *Output) PrintError(a ...any) {
	rnbw.ForgroundColor(rnbw.Red)
	o.Print("[ERROR] ")
	o.Print(a...)
	rnbw.ResetColor()
}

// PrintlnError prints an error message followed by a newline.
func (o *Output) PrintlnError(a ...any) {
	rnbw.ForgroundColor(rnbw.Red)
	o.Print("[ERROR] ")
	o.Println(a...)
	rnbw.ResetColor()
}

// PrintlnError prints an error message followed by a newline.
func (o *Output) PrintfError(format string, a ...any) {
	rnbw.ForgroundColor(rnbw.Red)
	o.Print("[ERROR] ")
	o.Printf(format, a...)
	rnbw.ResetColor()
}
