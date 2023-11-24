package errs

import "fmt"

type VisitorCodeInvalid struct {
	KameleoonError
}

func NewVisitorCodeInvalid(visitorCode string) *VisitorCodeInvalid {
	msg := fmt.Sprintf("Visitor code '%s' is not valid", visitorCode)
	return &VisitorCodeInvalid{NewKameleoonError(msg)}
}
