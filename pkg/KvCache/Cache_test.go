package KvCache

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var k1 = "11"
var k2 = 11
var v = 1
var v2 = 2

func Test_Cache_PutNew(t *testing.T) {
	cache := MakeCache[string, int]()
	cache.Put("11", 1)
	assert.Equal(t, 1, len(cache.values))
	for k, v := range cache.values {
		assert.Equal(t, "11", k)
		assert.Equal(t, 1, v)
	}
}

func Test_Cache_Get(t *testing.T) {
	cache := MakeCache[string, int]()
	_, found := cache.Get("11")
	assert.Equal(t, false, found)
	cache.Put("11", 1)
	v, found := cache.Get("11")
	assert.Equal(t, true, found)
	assert.Equal(t, 1, v)
}

func Test_Cache_PutOverwrite(t *testing.T) {
	cache := MakeCache[string, int]()
	cache.Put("11", 1)
	cache.Put("11", 2)
	v, found := cache.Get("11")
	assert.Equal(t, true, found)
	assert.Equal(t, 2, v)
}

func Test_Cache2_PutNew(t *testing.T) {
	cache := MakeCache2[string, int, int]()
	cache.Put("11", 11, 1)
	assert.Equal(t, 1, len(cache.values))
	for k1, l2values := range cache.values {
		assert.Equal(t, k1, "11")
		assert.Equal(t, 1, len(l2values))
		for k2, v := range l2values {
			assert.Equal(t, 11, k2)
			assert.Equal(t, 1, v)
		}
	}
}

func Test_Cache2_Get(t *testing.T) {
	cache := MakeCache2[string, int, int]()
	_, found := cache.Get(k1, k2)
	assert.Equal(t, false, found)
	cache.Put(k1, k2, v)
	gotV, found := cache.Get(k1, k2)
	assert.Equal(t, true, found)
	assert.Equal(t, v, gotV)
}

func Test_Cache2_PutOverwrite(t *testing.T) {
	cache := MakeCache2[string, int, int]()
	cache.Put(k1, k2, v)
	cache.Put(k1, k2, v2)
	gotV, found := cache.Get(k1, k2)
	assert.Equal(t, true, found)
	assert.Equal(t, v2, gotV)
}
