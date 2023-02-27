package kameleoon

import (
	"log"
	"os"
)

type Logger interface {
	Printf(format string, v ...interface{})
}

var defaultLogger Logger = log.New(os.Stdout, "KameleoonClient SDK", log.LstdFlags)
