package errs

import "fmt"

type VisitorCodeInvalid struct {
	KameleoonError
}

func NewVisitorCodeInvalid(errorDescription string) *VisitorCodeInvalid {
	msg := fmt.Sprintf("Visitor code %s", errorDescription)
	return &VisitorCodeInvalid{NewKameleoonError(msg)}
}
