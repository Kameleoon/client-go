package errs

// KameleeonError is the base type for kameleoon errors.
type KameleoonError struct {
	Message string `json:"message"`
}

func NewKameleoonError(msg string) KameleoonError {
	return KameleoonError{Message: "KameleoonClient SDK: " + msg}
}

func (e *KameleoonError) Error() string {
	return e.Message
}
