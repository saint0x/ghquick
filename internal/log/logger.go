package log

import (
	"fmt"
	"os"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
)

// Logger provides pretty console logging
type Logger struct {
	debug bool
}

// New creates a new logger instance
func New(debug bool) *Logger {
	return &Logger{debug: debug}
}

// Info prints an info message with a blue info icon
func (l *Logger) Info(format string, args ...interface{}) {
	fmt.Printf("%sℹ️  INFO: %s%s\n", colorBlue, fmt.Sprintf(format, args...), colorReset)
}

// Success prints a success message with a green checkmark
func (l *Logger) Success(format string, args ...interface{}) {
	fmt.Printf("%s✅ SUCCESS: %s%s\n", colorGreen, fmt.Sprintf(format, args...), colorReset)
}

// Error prints an error message with a red X
func (l *Logger) Error(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s❌ ERROR: %s%s\n", colorRed, fmt.Sprintf(format, args...), colorReset)
}

// Warning prints a warning message with a yellow warning icon
func (l *Logger) Warning(format string, args ...interface{}) {
	fmt.Printf("%s⚠️  WARNING: %s%s\n", colorYellow, fmt.Sprintf(format, args...), colorReset)
}

// Debug prints a debug message if debug mode is enabled
func (l *Logger) Debug(format string, args ...interface{}) {
	if l.debug {
		fmt.Printf("%s🔍 DEBUG: %s%s\n", colorPurple, fmt.Sprintf(format, args...), colorReset)
	}
}

// Step prints a step message with a cyan arrow
func (l *Logger) Step(format string, args ...interface{}) {
	fmt.Printf("%s➡️  %s%s\n", colorCyan, fmt.Sprintf(format, args...), colorReset)
}

// Command prints a command that's being executed
func (l *Logger) Command(cmd string, args ...string) {
	if l.debug {
		fullCmd := fmt.Sprintf("%s %s", cmd, strings.Join(args, " "))
		fmt.Printf("%s$ %s%s\n", colorPurple, fullCmd, colorReset)
	}
}
