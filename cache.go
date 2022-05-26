package memorycache

import (
	"sync"
	"time"
)

type (
	cacheItem struct {
		Value      interface{}
		created    time.Time
		expireTime time.Duration
	}

	CacheOption struct {
		// 最大缓存数, 默认为100条
		MaxItems int
		// 最长缓存时间， 默认为24小时
		Expire time.Duration
	}

	MemoryCache struct {
		data   map[string]cacheItem
		keys   []string
		option CacheOption
		sync.RWMutex
	}
)

func InitCache(cp CacheOption) *MemoryCache {
	if cp.MaxItems == 0 {
		cp.MaxItems = 1000
	}
	if cp.Expire == 0 {
		cp.Expire = 24 * time.Hour
	}
	cache := &MemoryCache{option: cp}
	go cache.runClear()
	return cache
}

func (cache *MemoryCache) runClear() {
	for {
		func() {
			defer time.Sleep(1 * time.Second)
			if len(cache.keys) > cache.option.MaxItems {
				key := cache.keys[len(cache.keys)-1]
				cache.Delete(key)
			}
			for _, key := range cache.keys {
				cache.RLock()
				it, ok := cache.data[key]
				cache.RUnlock()
				if ok {
					if t := time.Since(it.created); t > it.expireTime || t > cache.option.Expire {
						cache.Delete(key)
					}
				} else {
					cache.Lock()
					deleteString(cache.keys, key)
					cache.Unlock()
				}
			}
		}()
	}
}

// expire 过期时间，当 expire = 0时数据常驻内存，不超过最大最大缓存数或不超过最长缓存时间时不会过期
func (cache *MemoryCache) Add(key string, val interface{}, expire time.Duration) {
	cache.Lock()
	defer cache.Unlock()
	if cache.data == nil {
		cache.data = make(map[string]cacheItem)
	}
	it := cacheItem{
		Value:      val,
		created:    time.Now(),
		expireTime: expire,
	}
	cache.data[key] = it
	cache.keys = append(cache.keys, key)
}

func (cache *MemoryCache) Get(key string) (val interface{}, ok bool) {
	cache.RLock()
	defer cache.RUnlock()
	it, ok := cache.data[key]
	if ok {
		val = it.Value
	}
	return
}

func (cache *MemoryCache) Delete(key string) {
	cache.Lock()
	defer cache.Unlock()
	if _, ok := cache.data[key]; ok {
		delete(cache.data, key)
		deleteString(cache.keys, key)
	}
}

func deleteString(r []string, s string) []string {
	var res = []string{}
	for _, x := range r {
		if x != s {
			res = append(res, x)
		}
	}
	return res
}
