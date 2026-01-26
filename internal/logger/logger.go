// Package logger provides a simple, singleton-based logging utility.
// It supports logging messages to a file with different severity levels (DEBUG, INFO, WARNING, ERROR)
// and includes timestamps. The logger uses a singleton pattern to ensure only one file handle is used,
// with thread-safe operations via mutexes. It performs synchronous logging by default but offers
// a Sync function for immediate file flushing if needed.
package logger

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// LogLevel represents the severity level of a log entry.
type LogLevel int

const (
	// DEBUG level for fine-grained informational events useful for debugging.
	DEBUG LogLevel = iota
	// INFO level for general informational messages.
	INFO
	// WARNING level for potentially harmful situations or rare errors.
	WARNING
	// ERROR level for error conditions that might still allow the application to continue running.
	ERROR
)

var (
	// instance holds the single logWriter instance.
	instance *logWriter
	// once ensures that Init is called only once, implementing the singleton pattern.
	once sync.Once
)

// logWriter encapsulates the file handle and its own mutex for thread-safe writing.
// This struct is not exported to enforce the singleton pattern.
type logWriter struct {
	file *os.File
	mu   sync.Mutex
}

// Init initializes the logger with the specified file path.
// It uses the sync.Once mechanism to ensure the logger is initialized only once throughout the application lifecycle.
// The file is opened in append mode with permissions 0644 (readable by owner and group, writable by owner).
// If initialization fails, subsequent Log calls will silently ignore writes.
// Note: This function should be called early in the application's startup process.
func Init(filepath string) error {
	var err error
	once.Do(func() {
		var file *os.File
		file, err = os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return
		}
		instance = &logWriter{file: file}
	})
	return err
}

// Log writes a formatted log message to the initialized log file.
// It prepends the current timestamp and log level to the message.
func Log(level LogLevel, format string, a ...any) {
	w, err := getLogWriter()
	if err != nil {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	// Attempt to write to the file, ignoring potential errors for simplicity.
	if w.file != nil {
		// Generate a timestamp
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		// Format the full log entry: [timestamp] [LEVEL] message
		fmt.Fprintf(w.file, "[%s] [%s] %s\n", timestamp, levelString(level), fmt.Sprintf(format, a...))
	}
}

// Close closes the log file if it's open.
// This should typically be called during application shutdown to release resources.
// Returns an error if closing the file fails; otherwise, nil.
func Close() error {
	if instance != nil && instance.file != nil {
		return instance.file.Close()
	}
	return nil
}

// Sync forces any buffered data in the log file to be written to disk immediately.
// This is useful for ensuring logs are persisted in real-time, such as in logging-critical applications.
// Returns an error if syncing fails or if the logger isn't initialized.
func Sync() error {
	w, err := getLogWriter()
	if err != nil {
		return err
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.file.Sync()
}

// getLogWriter retrieves the singleton logWriter instance.
func getLogWriter() (*logWriter, error) {
	if instance == nil {
		return nil, fmt.Errorf("Init() must be called first")
	}
	return instance, nil
}

// levelString converts a LogLevel to its string representation.
func levelString(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}
