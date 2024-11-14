package utils

func Deopt[T any](field interface{}, defaultValue T) (T, bool) {
	if v, ok := field.(T); ok {
		return v, true
	}
	return defaultValue, field == nil
}
