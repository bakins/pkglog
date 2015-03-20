package pkglog

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type (
	Logger struct {
		mu       sync.Mutex
		out      io.Writer
		packages map[string]Level
		cache    map[string]Level
		level    Level // default log level
		output   Outputter
	}

	// perhaps wrap each log entry in a struct with pkgpath level, line number, buffer, etc
	// easier to pass around a struct??

	// Level is the log level
	Level uint8
)

const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `os.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
)

func (level Level) String() string {
	switch level {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warning"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	case PanicLevel:
		return "panic"
	}

	return "unknown"
}

// ParseLevel takes a string level and returns the log level constant.
func ParseLevel(lvl string) (Level, error) {
	switch lvl {
	case "panic":
		return PanicLevel, nil
	case "fatal":
		return FatalLevel, nil
	case "error":
		return ErrorLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	}

	var l Level
	return l, fmt.Errorf("not a valid log Level: %q", lvl)
}

// New creates a new Logger.
func New(out io.Writer) *Logger {
	return &Logger{
		out:      out,
		packages: make(map[string]Level),
		cache:    make(map[string]Level),
		level:    WarnLevel,
		output:   &DefaultOutputter{},
	}
}

var std = New(os.Stderr)

// StandardLogger returns the global Logger.
func StandardLogger() *Logger {
	return std
}

// SetLogLevel sets the default log level.
func SetLogLevel(level Level) {
	std.SetLogLevel(level)
}

// SetPackageLogLevel sets the log level for a given package
func (l *Logger) SetPackageLogLevel(pkg string, level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.packages[pkg] = level
}

// SetPackageLogLevel sets the log level for a given package
func SetPackageLogLevel(pkg string, level Level) {
	std.SetPackageLogLevel(pkg, level)
}

// SetOutputter sets the outputter to use.
func (l *Logger) SetOutputter(o Outputter) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.output = o
}

// SetOutputter sets the outputter to use.
func SetOutputter(o Outputter) {
	std.SetOutputter(o)
}

// LoggerWriter implements io.writer
type LoggerWriter struct {
	logger *Logger
}

// Write writes bytes to the logger at info level
func (w *LoggerWriter) Write(p []byte) (n int, err error) {
	// I have no clue..
	e := w.logger.newEntry(3)
	e.Output(e.Level, string(p))
	return len(p), nil
}

// Writer returns a struct which implements io.Writer
func (l *Logger) Writer() *LoggerWriter {
	return &LoggerWriter{logger: l}
}

// Printf logs a message at level Info.
func (l *Logger) Printf(format string, v ...interface{}) {
	e := l.newEntry(2)
	if e.Level >= InfoLevel {
		e.Output(InfoLevel, format, v...)
	}
}

// Printf logs a message at level Info.
func Printf(format string, v ...interface{}) {
	e := std.newEntry(2)
	if e.Level >= InfoLevel {
		e.Output(InfoLevel, format, v...)
	}
}

// Info logs a message at level Info.
func Info(format string, v ...interface{}) {
	e := std.newEntry(2)
	if e.Level >= InfoLevel {
		e.Output(InfoLevel, format, v...)
	}
}

// SetLogLevel sets the default log level.
func (l *Logger) SetLogLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}
