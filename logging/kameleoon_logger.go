package logging

import (
	"fmt"
)

type kameleoonLogger struct {
	logger   interface{}
	logLevel LogLevel
}

var kLogger = kameleoonLogger{
	logger:   NewDefaultLogger(),
	logLevel: WARNING,
}

func SetLogger(logger LoggerWithLevel) {
	kLogger.logger = logger
}

// DEPRECATED. Please use `logging.SetLogger(logging.LoggerWithLevel)` instead.
func SetOldLogger(logger Logger) {
	kLogger.logger = logger
}

func SetLogLevel(level LogLevel) {
	kLogger.logLevel = level
}

func GetLogLevel() LogLevel {
	return kLogger.logLevel
}

func Log(level LogLevel, data interface{}, args ...interface{}) {
	if checkLevel(level) {
		var message string
		switch v := data.(type) {
		case func() string:
			message = v()
		case string:
			if len(args) == 0 {
				message = v
			} else {
				message = fmt.Sprintf(v, prepareArgs(args...)...)
			}
		default:
			message = fmt.Sprintf("unsupported data type: %T", v)
		}
		writeMessage(level, message)
	}
}

func checkLevel(level LogLevel) bool {
	return level <= kLogger.logLevel && level != NONE
}

func writeMessage(level LogLevel, message string) {
	message = fmt.Sprintf("Kameleoon [%s]: %s", level, message)
	switch l := kLogger.logger.(type) {
	case Logger:
		l.Printf(message)
	case LoggerWithLevel:
		l.Log(level, message)
	}
}

func Info(data interface{}, args ...interface{}) {
	Log(INFO, data, args...)
}

func Error(data interface{}, args ...interface{}) {
	Log(ERROR, data, args...)
}

func Warning(data interface{}, args ...interface{}) {
	Log(WARNING, data, args...)
}

func Debug(data interface{}, args ...interface{}) {
	Log(DEBUG, data, args...)
}
