package network

import "fmt"

type ErrUnexpectedResponseStatus struct {
	Code int
}

func (err ErrUnexpectedResponseStatus) Error() string {
	return fmt.Sprintf("Received unexpected status code '%d'", err.Code)
}
