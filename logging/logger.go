package logging

import (
	"log"
	"os"
)

type LoggerMode int

const (
	Silent LoggerMode = iota
	Verbose
)

type Logger interface {
	Printf(format string, v ...interface{})
}

type LoggerImpl struct {
	mode   LoggerMode
	logger Logger
}

func NewLogger(mode LoggerMode, innerLogger Logger) Logger {
	return &LoggerImpl{
		mode:   mode,
		logger: innerLogger,
	}
}

func (l *LoggerImpl) Printf(format string, v ...interface{}) {
	if l.mode >= Verbose {
		l.logger.Printf(format, v...)
	}
}

var DefaultLogger Logger = log.New(os.Stdout, "KameleoonClient SDK: ", log.LstdFlags)
