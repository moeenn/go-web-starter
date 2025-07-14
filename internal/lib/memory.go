package lib

func Ref[T any](value T) *T {
	return &value
}
