package ex

func Map[T any, R any](collection []T, fn func(T) R) []R {
	if collection == nil {
		return nil
	}
	var res = make([]R, len(collection))
	for index, item := range collection {
		res[index] = fn(item)
	}
	return res
}

func Filtered[T any](collection []T, fn func(T) bool) []T {
	if collection == nil {
		return nil
	}
	var res = make([]T, 0)
	for _, item := range collection {
		if fn(item) {
			res = append(res, item)
		}
	}
	return res
}

func FindFirst[T any](collection []T, fn func(T) bool) (T, bool) {
	if collection == nil {
		var def T
		return def, false
	}
	for _, item := range collection {
		if fn(item) {
			return item, true
		}
	}
	var def T
	return def, false
}

func Contains[T comparable](collection []T, value T) bool {
	if collection == nil {
		return false
	}
	for _, item := range collection {
		if item == value {
			return true
		}
	}
	return false
}

func MapValues[M ~map[K]V, K comparable, V any](m M) []V {
	r := make([]V, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}
	return r
}
