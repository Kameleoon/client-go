package logging

type LoggerMode int

const (
	Silent LoggerMode = iota
	Verbose
)

type LogLevel int

const (
	NONE LogLevel = iota
	ERROR
	WARNING
	INFO
	DEBUG
)

var logLevelStrings = [...]string{"NONE", "ERROR", "WARNING", "INFO", "DEBUG"}

func (ll LogLevel) String() string {
	if int(ll) < 0 || int(ll) >= 5 {
		return "UNKNOWN"
	}
	return logLevelStrings[ll]
}

// DEPRECATED. Please use `logging.LoggerWithLevel` instead.
type Logger interface {
	Printf(format string, v ...interface{})
}

// DEPRECATED.
type LoggerImpl struct {
	mode   LoggerMode
	logger Logger
}

// DEPRECATED.
func NewLogger(mode LoggerMode, innerLogger Logger) Logger {
	return &LoggerImpl{
		mode:   mode,
		logger: innerLogger,
	}
}

// DEPRECATED.
func (l *LoggerImpl) Printf(format string, v ...interface{}) {
	if l.mode >= Verbose {
		l.logger.Printf(format, v...)
	}
}

type LoggerWithLevel interface {
	Log(level LogLevel, message string)
}
