package cache

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLRUCache_BasicOperations(t *testing.T) {

	cache := NewOrderCache(3)

	cache.Set("key1", "value1")

	val, found := cache.Get("key1")

	assert.True(t, found)

	assert.Equal(t, "value1", val)

	val, found = cache.Get("nonexistent")

	assert.False(t, found)

	assert.Nil(t, val)

}

func TestLRUCache_Capacity(t *testing.T) {

	cache := NewOrderCache(2)

	cache.Set("key1", "value1")

	cache.Set("key2", "value2")

	cache.Set("key3", "value3")

	val, found := cache.Get("key1")

	assert.False(t, found)

	assert.Nil(t, val)

	val, found = cache.Get("key2")

	assert.True(t, found)

	assert.Equal(t, "value2", val)

	val, found = cache.Get("key3")

	assert.True(t, found)

	assert.Equal(t, "value3", val)

}

func TestLRUCache_LRUEviction(t *testing.T) {

	cache := NewOrderCache(3)

	cache.Set("key1", "value1")

	cache.Set("key2", "value2")

	cache.Set("key3", "value3")

	cache.Get("key1")

	cache.Set("key4", "value4")

	val, found := cache.Get("key2")

	assert.False(t, found)

	assert.Nil(t, val)

	val, found = cache.Get("key1")

	assert.True(t, found)

	assert.Equal(t, "value1", val)

	val, found = cache.Get("key3")

	assert.True(t, found)

	assert.Equal(t, "value3", val)

	val, found = cache.Get("key4")

	assert.True(t, found)

	assert.Equal(t, "value4", val)

}

func TestLRUCache_UpdateExistingKey(t *testing.T) {

	cache := NewOrderCache(2)

	cache.Set("key1", "value1")

	cache.Set("key1", "updated_value")

	val, found := cache.Get("key1")

	assert.True(t, found)

	assert.Equal(t, "updated_value", val)

}

func TestLRUCache_ConcurrentAccess(t *testing.T) {

	cache := NewOrderCache(100)

	numGoroutines := 10

	iterations := 100

	var wg sync.WaitGroup

	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {

		go func(goroutineID int) {

			defer wg.Done()

			for j := 0; j < iterations; j++ {

				key := string(rune(goroutineID*1000 + j))

				cache.Set(key, key+"_value")

				cache.Get(key)

			}

		}(i)

	}

	wg.Wait()

	for i := 0; i < numGoroutines; i++ {

		for j := 0; j < iterations; j++ {

			key := string(rune(i*1000 + j))

			val, found := cache.Get(key)

			if found {

				assert.Equal(t, key+"_value", val)

			}

		}

	}

}

func TestLRUCache_EdgeCases(t *testing.T) {

	cache := NewOrderCache(0)

	cache.Set("key1", "value1")

	val, found := cache.Get("key1")

	assert.False(t, found)

	assert.Nil(t, val)

	cache = NewOrderCache(1)

	cache.Set("key1", "value1")

	cache.Set("key2", "value2")

	val, found = cache.Get("key1")

	assert.False(t, found)

	val, found = cache.Get("key2")

	assert.True(t, found)

	assert.Equal(t, "value2", val)

}

func TestLRUCache_OrderAfterAccess(t *testing.T) {

	cache := NewOrderCache(3)

	cache.Set("key1", "value1")

	cache.Set("key2", "value2")

	cache.Set("key3", "value3")

	cache.Get("key1")

	cache.Get("key3")

	cache.Set("key4", "value4")

	_, found := cache.Get("key2")

	assert.False(t, found)

	_, found = cache.Get("key1")

	assert.True(t, found)

	_, found = cache.Get("key3")

	assert.True(t, found)

	_, found = cache.Get("key4")

	assert.True(t, found)

}
