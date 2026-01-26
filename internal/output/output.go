// Package output supporting concurrent writes to stdout, files, or both.
// It provides a singleton pattern for output handling, with mode-based control.
package output

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/thxrsxm/rnbw"
)

// WriterType defines the targets where output can be written, allowing selective redirection.
type WriterType int

// Constants for WriterType to specify output destinations.
const (
	// Write to all available targets (stdout and file).
	ALL WriterType = iota
	// Write to stdout only.
	STDOUT
	// Write to file only.
	FILE
)

// out is a singleton struct for managing output writers and modes.
type out struct {
	// Map of writers keyed by type.
	writer map[WriterType]*outWriter
	// Current mode determining where output goes.
	mode WriterType
	// Mutex to ensure thread-safe operations on writers.
	mu sync.Mutex
}

var (
	// instance holds the single out instance.
	instance *out
	// once ensures that Init is called only once, implementing the singleton pattern.
	once sync.Once
)

// Init initializes the output system with a directory for output files and a flag for file writing.
// It sets up writers for stdout and optionally for a timestamped markdown file in the specified path.
func Init(outPath string, writeToFile bool) error {
	var err error
	once.Do(func() {
		instance = &out{writer: make(map[WriterType]*outWriter), mode: ALL}
		// Add stdout writer to output
		instance.writer[STDOUT] = newOutWriter(os.Stdout)
		// Prepare to write output to file
		if writeToFile {
			// Check out directory not exists
			if _, err := os.Stat(outPath); os.IsNotExist(err) {
				// Create out directory
				err := os.Mkdir(outPath, 0755)
				if err != nil {
					return
				}
			}
			// Set up out file with timestamp-based name
			outFilePath := fmt.Sprintf("%s/hzmind_%s.md", outPath, time.Now().Format("2006-01-02_15-04-05"))
			outFile, err := os.OpenFile(outFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				return
			}
			instance.writer[FILE] = newOutWriter(outFile)
		}
	})
	return err
}

// Close closes all output writers.
func Close() {
	if instance != nil {
		for _, v := range instance.writer {
			if closer, ok := v.writer.(io.Closer); ok {
				closer.Close()
			}
		}
	}
}

// Print prints the provided values to the current output targets without a newline.
func Print(a ...any) {
	o, err := getOut()
	if err != nil {
		return
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	for k, v := range o.writer {
		if o.mode != ALL && o.mode != k {
			continue
		}
		v.print(a...)
	}
}

// Println prints the provided values to the current output targets with a newline.
func Println(a ...any) {
	o, err := getOut()
	if err != nil {
		return
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	for k, v := range o.writer {
		if o.mode != ALL && o.mode != k {
			continue
		}
		v.println(a...)
	}
}

// Printf prints a formatted string to the current output targets.
func Printf(format string, a ...any) {
	o, err := getOut()
	if err != nil {
		return
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	for k, v := range o.writer {
		if o.mode != ALL && o.mode != k {
			continue
		}
		v.printf(format, a...)
	}
}

// PrintWarning prints a warning message with color styling to yellow.
func PrintWarning(a ...any) {
	rnbw.ForgroundColor(rnbw.Yellow)
	Print("[WARNING] ")
	Print(a...)
	rnbw.ResetColor()
}

// PrintlnWarning prints a warning message line with color styling.
func PrintlnWarning(a ...any) {
	rnbw.ForgroundColor(rnbw.Yellow)
	Print("[WARNING] ")
	Println(a...)
	rnbw.ResetColor()
}

// PrintfWarning prints a formatted warning message with color styling.
func PrintfWarning(format string, a ...any) {
	rnbw.ForgroundColor(rnbw.Yellow)
	Print("[WARNING] ")
	Printf(format, a...)
	rnbw.ResetColor()
}

// PrintError prints an error message with color styling to red.
func PrintError(a ...any) {
	rnbw.ForgroundColor(rnbw.Red)
	Print("[ERROR] ")
	Print(a...)
	rnbw.ResetColor()
}

// PrintlnError prints an error message line with color styling.
func PrintlnError(a ...any) {
	rnbw.ForgroundColor(rnbw.Red)
	Print("[ERROR] ")
	Println(a...)
	rnbw.ResetColor()
}

// PrintfError prints a formatted error message with color styling.
func PrintfError(format string, a ...any) {
	rnbw.ForgroundColor(rnbw.Red)
	Print("[ERROR] ")
	Printf(format, a...)
	rnbw.ResetColor()
}

// SetWriteMode sets the output mode to control where writes occur.
func SetWriteMode(mode WriterType) {
	o, err := getOut()
	if err != nil {
		return
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	o.mode = mode
}

// getOut retrieves the singleton out instance, or returns an error if not initialized.
func getOut() (*out, error) {
	if instance == nil {
		return nil, fmt.Errorf("Init() must be called first")
	}
	return instance, nil
}
