package cache

import (
	"github.com/hashicorp/golang-lru"
)

//Lru object
type Lru struct {
	cache *lru.Cache
}

//NewLru cache
func NewLru() *Lru {
	l, err := lru.New(50000)
	if err != nil {
		panic("failed to init lru cache")
	}
	return &Lru{cache: l}
}

//Add key and value to lru
func (l *Lru) Add(key, value interface{}) bool {
	return l.cache.Add(key, value)
}

//Get return value by key
func (l *Lru) Get(key interface{}) (interface{}, bool) {
	return l.cache.Get(key)
}

//Len return cache size
func (l *Lru) Len() int {
	return l.cache.Len()
}
