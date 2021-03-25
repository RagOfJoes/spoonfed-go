package logger

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// Singleton object
var logger *Logger

func init() {
	logger = New()
}

// New initializes the standard logger
func New() *Logger {
	base := logrus.New()
	l := &Logger{base}
	l.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}

	return l
}

// Exposed log fn:

// Debug Log
func Debug(args ...interface{}) {
	logger.Debugln(args...)
}

// Debugf Log
func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

// Errorfn Log errors of a [fn] with format
func Errorfn(fn string, err error) error {
	outerr := fmt.Errorf("[%s]: %v", fn, err)
	logger.Errorln(outerr)
	return outerr
}

// Info Log
func Info(args ...interface{}) {
	logger.Infoln(args...)
}

// Infof Log
func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

// Warn Log
func Warn(args ...interface{}) {
	logger.Warnln(args...)
}

// Warnf Log
func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

// Panic Log
func Panic(args ...interface{}) {
	logger.Panicln(args...)
}

// Panicf Log
func Panicf(format string, args ...interface{}) {
	logger.Panicf(format, args...)
}

// Error Log
func Error(args ...interface{}) {
	logger.Errorln(args...)
}

// Errorf Log
func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

// Fatal Log
func Fatal(args ...interface{}) {
	logger.Fatalln(args...)
}

// Fatalf Log
func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}
