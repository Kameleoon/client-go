package logging

import (
	"log"
	"os"
)

type DefaultLogger struct {
	logger *log.Logger
}

func NewDefaultLogger() LoggerWithLevel {
	return &DefaultLogger{
		logger: log.New(os.Stdout, "Kameleoon SDK: ", log.LstdFlags),
	}
}

func (dl DefaultLogger) Log(level LogLevel, message string) {
	dl.logger.Println(message)
}
