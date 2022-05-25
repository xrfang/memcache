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

// expire 过期时间，当 expire = 0时数据常驻内存，不会过期
func Add(key string, val interface{}, expire time.Duration) {
	cache.Lock()
	defer cache.Unlock()
	if cache.data == nil {
		cache.data = make(map[string]cacheItem)
	}
	if it, ok := cache.data[key]; ok && it.timer != nil {
		it.timer.Stop()
	}
	it := cacheItem{
		Value:      val,
		created:    time.Now(),
		expireTime: expire,
	}
	if expire > 0 {
		it.timer = time.AfterFunc(expire, func() {
			Delete(key)
		})
	}
	cache.data[key] = it
}

func Get(key string) (val interface{}, ok bool) {
	cache.RLock()
	defer cache.RUnlock()
	if it, ok := cache.data[key]; ok {
		val = it.Value
	}
	return
}

func Delete(key string) {
	cache.Lock()
	defer cache.Unlock()
	if it, ok := cache.data[key]; ok {
		if it.timer != nil {
			it.timer.Stop()
		}
		delete(cache.data, key)
	}
}
