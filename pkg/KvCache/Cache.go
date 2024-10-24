package KvCache

import (
	"sync"
)

type Cache[K comparable, V any] struct {
	mutex  sync.RWMutex
	values map[K]V
}

func MakeCache[K comparable, V any]() Cache[K, V] {
	return Cache[K, V]{
		values: make(map[K]V),
		mutex:  sync.RWMutex{},
	}
}

func (cache *Cache[K, V]) Put(key K, value V) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	cache.values[key] = value
}

func (cache *Cache[K, V]) Get(key K) (V, bool) {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()
	value, found := cache.values[key]
	return value, found
}

type Cache2[K1 comparable, K2 comparable, V any] struct {
	mutex  sync.RWMutex
	values map[K1]map[K2]V
}

func MakeCache2[K1 comparable, K2 comparable, V any]() Cache2[K1, K2, V] {
	return Cache2[K1, K2, V]{
		values: make(map[K1]map[K2]V),
		mutex:  sync.RWMutex{},
	}
}

func (cache *Cache2[K1, K2, V]) Put(keyPart1 K1, keyPart2 K2, value V) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	l2values, found := cache.values[keyPart1]
	if !found {
		l2values = make(map[K2]V)
		cache.values[keyPart1] = l2values
	}
	l2values[keyPart2] = value
}

func (cache *Cache2[K1, K2, V]) Get(keyPart1 K1, keyPart2 K2) (V, bool) {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()
	l2values, found := cache.values[keyPart1]
	if !found {
		var v V
		return v, false
	}
	value, found := l2values[keyPart2]
	return value, found
}
