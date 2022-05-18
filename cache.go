package memorycache

import (
	"sync"
	"time"
)

type (
	cacheItem struct {
		Key        string
		Value      interface{}
		created    time.Time
		expireTime time.Duration
		timer      *time.Timer
	}

	memoryCache struct {
		data map[string]cacheItem
		sync.RWMutex
	}
)

var cache *memoryCache

func init() {
	cache = &memoryCache{}
}

func Add(key string, val interface{}, expire time.Duration) {
	cache.Lock()
	defer cache.Unlock()
	if cache.data == nil {
		cache.data = make(map[string]cacheItem)
	}
	if it, ok := cache.data[key]; ok {
		it.timer.Stop()
	}
	cache.data[key] = cacheItem{
		Value:      val,
		created:    time.Now(),
		expireTime: expire,
		timer: time.AfterFunc(expire, func() {
			Delete(key)
		}),
	}
}

func Get(key string) (val interface{}, ok bool) {
	cache.RLock()
	defer cache.RUnlock()
	val, ok = cache.data[key]
	return
}

func Delete(key string) {
	cache.Lock()
	defer cache.Unlock()
	if it, ok := cache.data[key]; ok {
		it.timer.Stop()
		delete(cache.data, key)
	}
}
