package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Lru(t *testing.T) {
	lru := NewLru()
	lru.Add("test", []string{"aaa", "bbb"})
	val, ok := lru.Get("test")
	assert.True(t, ok)
	assert.Equal(t, val.([]string)[0], "aaa")
	assert.Equal(t, val.([]string)[1], "bbb")
	assert.Equal(t, lru.Len(), 1)
}
