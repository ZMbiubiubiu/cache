package lru

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type valueEntry string

func (v valueEntry) Len() int { return len(v) }

func TestCacheBasic(t *testing.T) {
	c := NewCache(100_0000, nil)
	key, value := "hello", valueEntry("world")
	err := c.Set(key, value)
	assert.Nil(t, err)

	got, ok := c.Get(key)
	assert.Equal(t, true, ok)
	assert.Equal(t, value, got)

	// delete and get again
	c.Delete(key)
	_, ok = c.Get(key)
	assert.Equal(t, false, ok)
}

func TestCacheLRUPolicy(t *testing.T) {
	c := NewCache(10, nil)
	key1, key2, key3 := "1", "2", "3"
	value1, value2, value3 := valueEntry("1000"), valueEntry("2000"), valueEntry("3000")

	// set and get k1/v1
	err := c.Set(key1, value1)
	assert.Nil(t, err)
	got1, ok := c.Get(key1)
	assert.Equal(t, true, ok)
	assert.Equal(t, value1, got1)

	// set and get k2/v2
	err = c.Set(key2, value2)
	assert.Nil(t, err)
	got2, ok := c.Get(key2)
	assert.Equal(t, true, ok)
	assert.Equal(t, value2, got2)

	// set and get k2/v2
	err = c.Set(key3, value3)
	assert.Nil(t, err)
	got3, ok := c.Get(key3)
	assert.Equal(t, true, ok)
	assert.Equal(t, value3, got3)
	// meanwhile, k2/v2 should exist but k1/v1 should be gone
	got2, ok = c.Get(key2)
	assert.Equal(t, true, ok)
	assert.Equal(t, value2, got2)

	_, ok = c.Get(key1)
	assert.Equal(t, false, ok)
}
