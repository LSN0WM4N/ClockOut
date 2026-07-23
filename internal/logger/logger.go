package logger

import (
	"fmt"
	"os"
	"strings"

	"ClockOut/internal/utils"
)

type level string

const (
	levelInfo  level = "INFO "
	levelWarn  level = "WARN "
	levelError level = "ERROR"
	levelDebug level = "DEBUG"
	levelFatal level = "FATAL"
)

// isDebugEnabled checks whether the DEBUG env var is truthy.
func isDebugEnabled() bool {
	v := strings.ToLower(os.Getenv("DEBUG"))
	return v == "true" || v == "1"
}

// output writes a formatted log line to the given writer.
func output(w *os.File, lvl level, context string, args ...string) {
	msg := strings.Join(args, " ")
	fmt.Fprintf(w, "[%s][%s][%s] %s.\n", lvl, utils.Now(), strings.ToUpper(context), msg)
}

// Print logs an informational message to STDOUT.
func Print(context string, args ...string) {
	output(os.Stdout, levelInfo, context, args...)
}

// Warn logs a warning message to STDOUT.
func Warn(context string, args ...string) {
	output(os.Stdout, levelWarn, context, args...)
}

// Error logs an error message to STDERR.
func Error(context string, args ...string) {
	output(os.Stderr, levelError, context, args...)
}

// Debug logs a debug message to STDOUT, only when DEBUG=true or DEBUG=1.
func Debug(context string, args ...string) {
	if isDebugEnabled() {
		output(os.Stdout, levelDebug, context, args...)
	}
}

// Fatal logs the error with context to STDERR and exits with status 1.
func Fatal(context string, err error) {
	output(os.Stderr, levelFatal, context, err.Error())
	os.Exit(1)
}
