// Package input provides a singleton-based input handler for reading user input and passwords securely.
// It supports thread-safe operations with mutexes to handle concurrent reads.
package input

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"syscall"

	"github.com/thxrsxm/harzmind-code/internal/output"
	"golang.org/x/term"
)

// in represents a singleton input handler with a buffered reader and a mutex for thread safety.
type in struct {
	// Buffered reader for reading input from stdin.
	reader *bufio.Reader
	// Mutex to ensure thread-safe operations on writers.
	mu sync.Mutex
}

var (
	// instance holds the single in instance.
	instance *in
	// once ensures that Init is called only once, implementing the singleton pattern.
	once sync.Once
)

// Init initializes the input handler by creating a singleton instance with a buffered reader for os.Stdin.
// It ensures the instance is created only once across the application's lifecycle.
func Init() error {
	var err error
	once.Do(func() {
		instance = &in{reader: bufio.NewReader(os.Stdin)}
	})
	return err
}

// ReadInput reads a line of input from the user.
// It also writes the input to the output file if available.
// Returns the trimmed input string and any error encountered.
func ReadInput(writeToFile bool) (string, error) {
	i, err := getIn()
	if err != nil {
		return "", err
	}
	i.mu.Lock()
	defer i.mu.Unlock()
	input, err := i.reader.ReadString('\n')
	// Write user input to output file
	if writeToFile {
		output.SetWriteMode(output.FILE)
		output.Print(input)
		output.SetWriteMode(output.ALL)
	}
	// Handle input error
	if err != nil {
		return "", err
	}
	input = strings.TrimSpace(input)
	return input, nil
}

// ReadPassword reads a password from the user securely without echoing input to the terminal.
// It uses the term package to read from stdin and outputs a newline after reading.
// Returns the password string and any error encountered.
func ReadPassword() (string, error) {
	i, err := getIn()
	if err != nil {
		return "", err
	}
	i.mu.Lock()
	defer i.mu.Unlock()
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	output.Println()
	return string(bytePassword), nil
}

// getIn retrieves the singleton in instance, or returns an error if not initialized.
func getIn() (*in, error) {
	if instance == nil {
		return nil, fmt.Errorf("Init() must be called first")
	}
	return instance, nil
}
