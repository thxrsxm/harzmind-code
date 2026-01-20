// Package logger provides a simple file-based logging utility for the application,
// supporting different log levels (Debug, Info, Warning, Error) with timestamped entries.
// It appends log messages to a specified file and provides methods to log at various levels.
package logger

import (
	"fmt"
	"os"
	"time"
)

// LogLevel represents the severity level of a log entry.
type LogLevel int

const (
	Debug   LogLevel = iota // Debug level for fine-grained informational events useful for debugging.
	Info                    // Info level for general informational messages.
	Warning                 // Warning level for potentially harmful situations or rare errors.
	Error                   // Error level for error conditions that might still allow the application to continue running.
)

// Logger represents a file-based logger.
// It holds a reference to an open file where log entries are written.
type Logger struct {
	file *os.File
}

// NewLogger creates a new Logger instance that writes to the specified file path.
// It opens or creates the file in append mode with read-write permissions for the owner.
// Returns an error if the file cannot be opened or created.
// The caller is responsible for calling Close() to free resources.
func NewLogger(filePath string) (*Logger, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &Logger{file: file}, nil
}

// Close closes the log file if it is open.
// It returns nil if the file is already closed or does not exist.
// Should be called to ensure proper resource cleanup, typically deferred after NewLogger.
func (l *Logger) Close() error {
	if l.file == nil {
		return nil
	}
	return l.file.Close()
}

// log is a private helper method that formats and writes a log entry at the given level.
// It prepends a timestamp and log level to the message, then appends a newline.
// Ignores write errors to prevent logging failures from crashing the application.
func (l *Logger) log(level LogLevel, format string, a ...any) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	// Format the full log entry: [timestamp] [LEVEL] message
	logEntry := fmt.Sprintf("[%s] [%s] %s\n", timestamp, levelString(level), fmt.Sprintf(format, a...))
	// Attempt to write to the file, ignoring potential errors for simplicity.
	if l.file != nil {
		_, _ = l.file.WriteString(logEntry)
	}
}

// Debugf logs a message at the Debug level using fmt.Sprintf-style formatting.
// Useful for verbose debugging information that might not be needed in production.
func (l *Logger) Debugf(format string, a ...any) {
	l.log(Debug, format, a...)
}

// Infof logs a message at the Info level using fmt.Sprintf-style formatting.
// Suitable for general informational messages about application events.
func (l *Logger) Infof(format string, a ...any) {
	l.log(Info, format, a...)
}

// Warningf logs a message at the Warning level using fmt.Sprintf-style formatting.
// Indicates potential issues or warnings that do not halt execution.
func (l *Logger) Warningf(format string, a ...any) {
	l.log(Warning, format, a...)
}

// Errorf logs a message at the Error level using fmt.Sprintf-style formatting.
// Used for serious errors that may require attention but allow continued operation.
func (l *Logger) Errorf(format string, a ...any) {
	l.log(Error, format, a...)
}

// levelString converts a LogLevel to its string representation.
// Used internally for log entry formatting. Returns "UNKNOWN" for invalid levels.
func levelString(level LogLevel) string {
	switch level {
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warning:
		return "WARNING"
	case Error:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}
