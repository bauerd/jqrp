package log

import (
	"fmt"
	"log"
	"os"
)

// Level is a log level.
type Level uint

func (l Level) String() string {
	switch l {
	case Debug:
		return "DBG"
	case Info:
		return "INF"
	default:
		return "ERR"
	}
}

const (
	// Debug is the lowest log level leading to most output.
	Debug Level = iota

	// Info is the log level for general-purpose messages.
	Info

	// Error is the highest log level reserved for error conditions.
	Error
)

// Logger is a levelled logger.
type Logger struct {
	logger *log.Logger
	Level
}

// New returns a new logger with given level that writes to stdout.
func New(level Level) *Logger {
	return &Logger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
		Level:  level,
	}
}

// Debug logs a debug message.
func (l *Logger) Debug(msg string) {
	l.log(Debug, msg)
}

// Info logs an info message.
func (l *Logger) Info(msg string) {
	l.log(Info, msg)
}

// Error logs an error message.
func (l *Logger) Error(msg string) {
	l.log(Error, msg)
}

func (l *Logger) log(level Level, msg string) {
	if level < l.Level {
		return
	}
	l.logger.Printf(fmt.Sprintf("[%s] %s", level, msg))
}
