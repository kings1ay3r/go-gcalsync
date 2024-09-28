package logger

import (
	"context"
	"fmt"
	"gorm.io/gorm/logger"
	"os"
	"sync"
	"time"

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
	return instance
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

// Trace TODO: Implement Tracer. Currently implementing structure to satisfy interface
func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
}

// LogMode TODO: Implement this. Currently implementing wrapper to satisfy interface
func (l *Logger) LogMode(level logger.LogLevel) logger.Interface {
	return NewLogger()
}
