package utils

import (
	"fmt"

	"github.com/Kameleoon/client-go/v3/errs"
)

const (
	VisitorCodeMaxLength = 255
	VisitorCodeLength    = 16
)

func ValidateVisitorCode(visitorCode string) error {
	if visitorCode == "" {
		return errs.NewVisitorCodeInvalid("is empty")
	} else if len(visitorCode) > VisitorCodeMaxLength {
		msg := fmt.Sprintf("is longer than %d chars", VisitorCodeMaxLength)
		return errs.NewVisitorCodeInvalid(msg)
	}
	return nil
}

func GenerateVisitorCode() string {
	return GetRandomString(VisitorCodeLength)
}
