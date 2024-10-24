package utils

func Flat[T any](slices [][]T) []T {
	var result = make([]T, 0)
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return result
}

func SliceToMap[T any, K comparable, V any](
	slice []T,
	key func(T) K,
	value func(T) V,
) map[K]V {
	result := make(map[K]V)
	for _, item := range slice {
		result[key(item)] = value(item)
	}
	return result
}

func MapSlice[T any, R any](src []T, fn func(T) R) []R {
	var result = make([]R, 0)
	for _, item := range src {
		result = append(result, fn(item))
	}
	return result
}

func FilterSlice[T any](slice []T, fn func(T) bool) []T {
	if slice == nil {
		return nil
	}
	var res = make([]T, 0)
	for _, item := range slice {
		if fn(item) {
			res = append(res, item)
		}
	}
	return res
}
