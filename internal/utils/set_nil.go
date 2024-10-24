package utils

func SetNil[T any](ptr **T) {
	*ptr = nil
}
