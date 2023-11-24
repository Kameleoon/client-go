package errs

type InternalError struct {
	KameleoonError
}

func NewInternalError(msg string) *InternalError {
	return &InternalError{NewKameleoonError(msg)}
}
