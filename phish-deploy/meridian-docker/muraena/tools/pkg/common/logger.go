package common

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fatih/color"
)

// Logger provides structured logging with color support
type Logger struct {
	verbose bool
	logFile *os.File
}

// NewLogger creates a new logger instance
func NewLogger(verbose bool, logPath string) (*Logger, error) {
	var logFile *os.File
	var err error

	if logPath != "" {
		logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
	}

	return &Logger{
		verbose: verbose,
		logFile: logFile,
	}, nil
}

// Info logs an informational message
func (l *Logger) Info(msg string) {
	l.log("INFO", color.CyanString("ℹ"), msg)
}

// Success logs a success message
func (l *Logger) Success(msg string) {
	l.log("SUCCESS", color.GreenString("✓"), msg)
}

// Warning logs a warning message
func (l *Logger) Warning(msg string) {
	l.log("WARNING", color.YellowString("⚠"), msg)
}

// Error logs an error message
func (l *Logger) Error(msg string) {
	l.log("ERROR", color.RedString("✗"), msg)
}

// Debug logs a debug message (only if verbose is enabled)
func (l *Logger) Debug(msg string) {
	if l.verbose {
		l.log("DEBUG", color.MagentaString("⚙"), msg)
	}
}

// log is the internal logging function
func (l *Logger) log(level, icon, msg string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMsg := fmt.Sprintf("[%s] %s %s", timestamp, icon, msg)

	fmt.Println(logMsg)

	if l.logFile != nil {
		log.SetOutput(l.logFile)
		log.Printf("[%s] %s\n", level, msg)
	}
}

// Close closes the log file if open
func (l *Logger) Close() {
	if l.logFile != nil {
		l.logFile.Close()
	}
}

// Banner prints a formatted banner
func Banner(title string) {
	width := 65
	border := "═"

	topBorder := "╔" + repeatString(border, width-2) + "╗"
	bottomBorder := "╚" + repeatString(border, width-2) + "╝"

	padding := (width - len(title) - 2) / 2
	titleLine := "║" + repeatString(" ", padding) + title + repeatString(" ", width-len(title)-padding-2) + "║"

	cyan := color.New(color.FgCyan).SprintFunc()
	fmt.Println(cyan(topBorder))
	fmt.Println(cyan(titleLine))
	fmt.Println(cyan(bottomBorder))
}

func repeatString(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
