package utils

func GroupBy[T any, K comparable, V any](
	slice []T,
	key func(T) K,
	value func(T) V,
) map[K][]V {
	result := make(map[K][]V)
	GroupByTo(slice, result, key, value)
	return result
}

func GroupByTo[T any, K comparable, V any](
	slice []T,
	dest map[K][]V,
	key func(T) K,
	value func(T) V,
) {
	for _, item := range slice {
		k := key(item)
		v := value(item)
		dest[k] = append(dest[k], v)
	}
}

func GroupItemTo[T any, K comparable, V any](
	item T,
	dest map[K][]V,
	key func(T) K,
	value func(T) V,
) {
	k := key(item)
	v := value(item)
	dest[k] = append(dest[k], v)
}

func CloneMap[K comparable, V any](src map[K]V) map[K]V {
	result := make(map[K]V)
	for k, v := range src {
		result[k] = v
	}
	return result
}

func TransformMapValues[K comparable, SRC_V any, DST_V any](
	srcMap map[K]SRC_V,
	transform func(key K, src SRC_V) DST_V,
) map[K]DST_V {
	result := make(map[K]DST_V)
	for key, src := range srcMap {
		result[key] = transform(key, src)
	}
	return result
}

func MapForEach[K comparable, V any](
	m map[K]V,
	block func(key K, value V),
) {
	for key, value := range m {
		block(key, value)
	}
}

func MapGetOrCreate[K comparable, V any](
	m map[K]V,
	key K,
	create func() V,
) V {
	result, found := m[key]
	if !found {
		result = create()
		m[key] = result
	}
	return result
}

func MapFilter[K comparable, V any](
	m map[K]V,
	condition func(key K, value V) bool,
) map[K]V {
	result := make(map[K]V)
	for key, value := range m {
		if condition(key, value) {
			result[key] = value
		}
	}
	return result
}
