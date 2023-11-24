package errs

type FeatureError struct {
	KameleoonError
}

func NewFeatureError(msg string) FeatureError {
	return FeatureError{NewKameleoonError("Feature Error: " + msg)}
}
