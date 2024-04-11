package logger

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	entry *logrus.Entry
}

func (e *Logger) Debug(args ...interface{}) {
	e.entry.Debug(args...)
}

func (e *Logger) Info(args ...interface{}) {
	e.entry.Info(args...)
}

func (e *Logger) Warning(args ...interface{}) {
	e.entry.Warning(args...)
}

func (e *Logger) Error(args ...interface{}) {
	e.entry.Error(args...)
}

func (e *Logger) Fatal(args ...interface{}) {
	e.entry.Fatal(args...)
}

func (e *Logger) Panic(args ...interface{}) {
	e.entry.Panic(args...)
}

func (e *Logger) Debugf(format string, args ...interface{}) {
	e.entry.Debugf(format, args...)
}

func (e *Logger) Infof(format string, args ...interface{}) {
	e.entry.Infof(format, args...)
}

func (e *Logger) Warningf(format string, args ...interface{}) {
	e.entry.Warningf(format, args...)
}

func (e *Logger) Errorf(format string, args ...interface{}) {
	e.entry.Errorf(format, args...)
}

func (e *Logger) Fatalf(format string, args ...interface{}) {
	e.entry.Fatalf(format, args...)
}

func (e *Logger) Panicf(format string, args ...interface{}) {
	e.entry.Panicf(format, args...)
}

var once sync.Once
var singleton *Logger

func logger() *Logger {
	once.Do(func() {
		entry := logrus.NewEntry(logrus.New())
		singleton = &Logger{
			entry: entry,
		}
	})

	return singleton
}

// ------ exported methods:
var (
	// FromContext 函数的别名
	F = FromContext
)

func FromContext(ctx context.Context) *Logger {
	srcLogger := logger()

	// ...extract info from ctx

	return srcLogger
}

func Debug(args ...interface{}) {
	logger().Debug(args...)
}

func Info(args ...interface{}) {
	logger().Info(args...)
}

func Warning(args ...interface{}) {
	logger().Warning(args...)
}

func Error(args ...interface{}) {
	logger().Error(args...)
}

func Fatal(args ...interface{}) {
	logger().Fatal(args...)
}

func Panic(args ...interface{}) {
	logger().Panic(args...)
}

func Debugf(format string, args ...interface{}) {
	logger().Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	logger().Infof(format, args...)
}

func Warningf(format string, args ...interface{}) {
	logger().Warningf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	logger().Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	logger().Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	logger().Panicf(format, args...)
}
