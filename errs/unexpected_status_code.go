package errs

import "fmt"

type UnexpectedStatusCode struct {
	KameleoonError
}

func NewUnexpectedStatusCode(code int, body []byte) *UnexpectedStatusCode {
	msg := fmt.Sprintf("Received unexpected status code: %d, body: %s", code, string(body[:]))
	return &UnexpectedStatusCode{NewKameleoonError(msg)}
}
