package kameleoon

import "errors"

var (
	ErrInvalidFeatureKeyType = errors.New("feature key should be a string or an int")
	ErrBadStatus             = errors.New("bad status code")
	ErrEmptyResponse         = errors.New("empty response")
)

// APIError is the base type for endpoint-specific errors.
type APIError struct {
	Message string `json:"message"`
}

func (e APIError) Error() string {
	return e.Message
}

func (e APIError) IsApiError() bool {
	return true
}

type ErrNotFound struct {
	APIError
}

func (e ErrNotFound) IsNotFoundError() bool {
	return true
}

func (e ErrNotFound) Error() string {
	return e.Message + " not found"
}

type ErrVariationNotFound struct {
	ErrNotFound
}

func newErrVariationNotFound(msg string) error {
	return &ErrVariationNotFound{ErrNotFound{APIError{Message: msg}}}
}

func (e ErrVariationNotFound) Error() string {
	return "variation " + e.ErrNotFound.Error()
}

type ErrExperimentConfigNotFound struct {
	ErrNotFound
}

func newErrExperimentConfigNotFound(msg string) error {
	return &ErrExperimentConfigNotFound{ErrNotFound{APIError{Message: msg}}}
}

func (e ErrExperimentConfigNotFound) Error() string {
	return "experiment " + e.ErrNotFound.Error()
}

type ErrFeatureConfigNotFound struct {
	ErrNotFound
}

func newErrFeatureConfigNotFound(msg string) error {
	return &ErrFeatureConfigNotFound{ErrNotFound{APIError{Message: msg}}}
}

func (e ErrFeatureConfigNotFound) Error() string {
	return "feature flag " + e.ErrNotFound.Error()
}

type ErrFeatureVariableNotFound struct {
	ErrNotFound
}

func newErrNotFound(msg string) error {
	return &ErrNotFound{APIError{Message: msg}}
}

func newErrFeatureVariableNotFound(msg string) error {
	return &ErrFeatureVariableNotFound{ErrNotFound{APIError{Message: msg}}}
}

func (e ErrFeatureVariableNotFound) Error() string {
	return "feature variable " + e.ErrNotFound.Error()
}

type ErrCredentialsNotFound struct {
	ErrNotFound
}

func newErrCredentialsNotFound(msg string) error {
	return &ErrCredentialsNotFound{ErrNotFound{APIError{Message: msg}}}
}

func (e ErrCredentialsNotFound) Error() string {
	return "credentials " + e.ErrNotFound.Error()
}

type ErrNotTargeted struct {
	APIError
}

func newErrNotTargeted(msg string) error {
	return &ErrNotTargeted{APIError{Message: msg}}
}

func (e ErrNotTargeted) Error() string {
	return "visitor " + e.Message + " is not targeted"
}

type ErrNotActivated struct {
	APIError
}

func newErrNotActivated(msg string) error {
	return &ErrNotActivated{APIError{Message: msg}}
}

func (e ErrNotActivated) Error() string {
	return "visitor " + e.Message + " is not activated"
}
