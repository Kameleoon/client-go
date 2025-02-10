package utils

func Reref[T any](ref *T) *T {
	if ref != nil {
		copy := new(T)
		*copy = *ref
		return copy
	}
	return nil
}
