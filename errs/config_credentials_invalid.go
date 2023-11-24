package errs

type ConfigCredentialsInvalid struct {
	ConfigError
}

func NewConfigCredentialsInvalid(msg string) *ConfigCredentialsInvalid {
	return &ConfigCredentialsInvalid{NewConfigError(msg)}
}
