package utils

func Must[T any](content T, err error) T {
	if err != nil {
		panic(err)
	}
	return content
}
