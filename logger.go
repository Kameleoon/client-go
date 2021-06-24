package kameleoon

import (
	"log"
	"os"
)

type Logger interface {
	Printf(format string, args ...interface{})
}

var defaultLogger Logger = log.New(os.Stdout, "", log.LstdFlags)
