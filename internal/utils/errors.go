package utils

func OnlyError[T any](value T, err error) error {
	return err
}
