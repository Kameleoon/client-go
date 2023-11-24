package errs

type ConfigError struct {
	KameleoonError
}

func NewConfigError(msg string) ConfigError {
	return ConfigError{NewKameleoonError("Config Error: " + msg)}
}
