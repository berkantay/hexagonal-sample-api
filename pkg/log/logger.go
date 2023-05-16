package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger provides a simple interface for logging messages using the logrus library.
type Logger struct {
	logger *logrus.Logger
}

// NewLogger creates a new instance of Logger that logs messages to a file.
func NewLogger(filename string) *Logger {
	logFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	logger := logrus.New()
	logger.SetOutput(logFile)
	return &Logger{
		logger: logger,
	}
}
func NewTestLogger() *Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	return &Logger{
		logger: logger,
	}
}

// Info logs an informational message to the logger with the specified message and optional data fields.
func (l *Logger) Info(message string, data ...map[string]interface{}) {
	logFields := logrus.Fields{}
	for _, dataMap := range data {
		for key, value := range dataMap {
			logFields[key] = value
		}
	}
	l.logger.WithFields(logFields).Info(message)
}

// Warn logs an warning message to the logger with the specified message and optional data fields.
func (l *Logger) Warn(message string, data ...map[string]interface{}) {
	logFields := logrus.Fields{}
	for _, dataMap := range data {
		for key, value := range dataMap {
			logFields[key] = value
		}
	}
	l.logger.WithFields(logFields).Warn(message)
}

// Fatal logs an fatal message to the logger with the specified message and optional data fields.
func (l *Logger) Fatal(message string, data ...map[string]interface{}) {
	logFields := logrus.Fields{}
	for _, dataMap := range data {
		for key, value := range dataMap {
			logFields[key] = value
		}
	}
	l.logger.WithFields(logFields).Fatal(message)
}

// Shuts down writer.
func (l *Logger) Close() {
	l.logger.Writer().Close()
}
