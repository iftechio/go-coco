package xecho

import (
	"io"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	xLogger "github.com/xieziyu/go-coco/utils/logger"
)

type logger struct{}

func (l *logger) Output() io.Writer                         { return os.Stdout }
func (l *logger) SetOutput(w io.Writer)                     {}
func (l *logger) Prefix() string                            { return "" }
func (l *logger) SetPrefix(p string)                        {}
func (l *logger) Level() log.Lvl                            { return 0 }
func (l *logger) SetLevel(v log.Lvl)                        {}
func (l *logger) SetHeader(h string)                        {}
func (l *logger) Debug(i ...interface{})                    { xLogger.Debug(i...) }
func (l *logger) Debugj(j log.JSON)                         { xLogger.Debug(j) }
func (l *logger) Debugf(format string, args ...interface{}) { xLogger.Debugf(format, args...) }
func (l *logger) Info(i ...interface{})                     { xLogger.Info(i...) }
func (l *logger) Infoj(j log.JSON)                          { xLogger.Info(j) }
func (l *logger) Infof(format string, args ...interface{})  { xLogger.Infof(format, args...) }
func (l *logger) Print(i ...interface{})                    { xLogger.Info(i...) }
func (l *logger) Printj(j log.JSON)                         { xLogger.Info(j) }
func (l *logger) Printf(format string, args ...interface{}) { xLogger.Infof(format, args...) }
func (l *logger) Warn(i ...interface{})                     { xLogger.Warning(i...) }
func (l *logger) Warnj(j log.JSON)                          { xLogger.Warning(j) }
func (l *logger) Warnf(format string, args ...interface{})  { xLogger.Warningf(format, args...) }
func (l *logger) Error(i ...interface{})                    { xLogger.Error(i...) }
func (l *logger) Errorj(j log.JSON)                         { xLogger.Error(j) }
func (l *logger) Errorf(format string, args ...interface{}) { xLogger.Errorf(format, args...) }
func (l *logger) Panic(i ...interface{})                    { xLogger.Panic(i...) }
func (l *logger) Panicj(j log.JSON)                         { xLogger.Panic(j) }
func (l *logger) Panicf(format string, args ...interface{}) { xLogger.Panicf(format, args...) }
func (l *logger) Fatal(i ...interface{})                    { xLogger.Fatal(i...) }
func (l *logger) Fatalj(j log.JSON)                         { xLogger.Fatal(j) }
func (l *logger) Fatalf(format string, args ...interface{}) { xLogger.Fatalf(format, args...) }

var _ echo.Logger = &logger{}
