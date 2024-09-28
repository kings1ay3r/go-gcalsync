package logger

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

// Logger ...
type Logger struct {
	*logrus.Logger
}

// singleton instance and sync.Once to ensure it's initialized only once
var (
	instance *Logger
	once     sync.Once
)

// GetInstance returns the singleton logger instance
func GetInstance() *Logger {
	once.Do(func() {
		instance = NewLogger() // Initialize the singleton instance using NewLogger
	})
	return instance // Return the singleton instance
}

// NewLogger ...
func NewLogger() *Logger {
	log := logrus.New()

	log.SetOutput(os.Stdout)

	log.SetLevel(logrus.InfoLevel)

	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: false,
	})

	return &Logger{log}
}

// Info logs an informational message
func (l *Logger) Info(ctx context.Context, msg string, args ...interface{}) {
	l.WithFields(logrus.Fields{}).Info(fmt.Sprintf(msg, args...))
	//l.WithFields(logrus.Fields{"context": ctx}).Info(fmt.Sprintf(msg, args...))
}

// Error logs an error message
func (l *Logger) Error(ctx context.Context, msg string, args ...interface{}) {
	l.WithFields(logrus.Fields{}).Error(fmt.Sprintf(msg, args...))
	//l.WithFields(logrus.Fields{"context": ctx}).Error(fmt.Sprintf(msg, args...))
}

// Warn logs a warning message
func (l *Logger) Warn(ctx context.Context, msg string, args ...interface{}) {
	l.WithFields(logrus.Fields{}).Warn(fmt.Sprintf(msg, args...))
	//l.WithFields(logrus.Fields{"context": ctx}).Warn(fmt.Sprintf(msg, args...))
}

// Debug logs a debug message
func (l *Logger) Debug(ctx context.Context, msg string, args ...interface{}) {
	l.WithFields(logrus.Fields{}).Debug(fmt.Sprintf(msg, args...))
	//l.WithFields(logrus.Fields{"context": ctx}).Debug(fmt.Sprintf(msg, args...))
}
