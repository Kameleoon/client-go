package errs

import "fmt"

type UnexpectedStatusCode struct {
	KameleoonError
}

func NewUnexpectedStatusCode(code int) *UnexpectedStatusCode {
	msg := fmt.Sprintf("Received unexpected status code '%d'", code)
	return &UnexpectedStatusCode{NewKameleoonError(msg)}
}
