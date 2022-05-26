package memorycache

import (
	"sort"
	"sync"
	"time"
)

type (
	cacheItem struct {
		data interface{}
		hits int64 //使用次数（LFU）
		used int64 //最后使用时间戳（毫秒，LRU）
		ttl  *time.Time
	}
	EvictionPolicy byte
	CacheOption    struct {
		Capacity int            //最大缓存数, 默认1024条
		Policy   EvictionPolicy //缓存清理策略，默认LRU
	}
	memoryCache struct {
		data map[string]*cacheItem
		opts *CacheOption
		wttl bool //包含有超时的条目
		sync.Mutex
	}
)

const (
	PolicyLRU EvictionPolicy = iota
	PolicyLFU
)

func New(co *CacheOption) *memoryCache {
	if co == nil {
		co = new(CacheOption)
	}
	if co.Capacity <= 0 {
		co.Capacity = 1024
	}
	mc := memoryCache{opts: co, data: make(map[string]*cacheItem)}
	go func() {
		for {
			time.Sleep(time.Second)
			mc.refresh()
			mc.evict()
		}
	}()
	return &mc
}

func (cache *memoryCache) refresh() {
	cache.Lock()
	defer cache.Unlock()
	if !cache.wttl {
		return
	}
	wttl := false
	for k, ci := range cache.data {
		if ci.ttl != nil {
			if time.Now().After(*ci.ttl) {
				delete(cache.data, k)
			} else {
				wttl = true
			}
		}
	}
	cache.wttl = wttl
}

func (cache *memoryCache) evict() {
	cache.Lock()
	defer cache.Unlock()
	over := len(cache.data) - cache.opts.Capacity
	if over <= 0 {
		return
	}
	type keyrank struct {
		key  string
		used int64
		hits int64
	}
	var krs []keyrank
	for k, ci := range cache.data {
		krs = append(krs, keyrank{key: k, used: ci.used, hits: ci.hits})
		switch cache.opts.Policy {
		case PolicyLRU:
			sort.Slice(krs, func(i, j int) bool {
				switch {
				case krs[i].used < krs[j].used:
					return true
				case krs[i].used > krs[j].used:
					return false
				default:
					switch {
					case krs[i].hits < krs[j].hits:
						return true
					case krs[i].hits > krs[j].hits:
						return false
					default:
						return krs[i].key < krs[j].key
					}
				}
			})
		default: //LFU
			sort.Slice(krs, func(i, j int) bool {
				switch {
				case krs[i].hits < krs[j].hits:
					return true
				case krs[i].hits > krs[j].hits:
					return false
				default:
					switch {
					case krs[i].used < krs[j].used:
						return true
					case krs[i].used > krs[j].used:
						return false
					default:
						return krs[i].key < krs[j].key
					}
				}
			})
		}
		if len(krs) > over {
			krs = krs[:over]
		}
	}
	for _, kr := range krs {
		delete(cache.data, kr.key)
	}
}

//expire为可选的过期时间，当expire>0时，即便缓存没有满，数据也会因超时被清理
func (cache *memoryCache) Set(key string, val any, expire ...time.Duration) {
	cache.Lock()
	defer cache.Unlock()
	ci := cache.data[key]
	if ci == nil {
		ci = new(cacheItem)
	}
	ci.data = val
	ci.used = time.Now().UnixMilli()
	ci.hits++
	if len(expire) > 0 && expire[0] > 0 {
		ttl := time.Now().Add(expire[0])
		ci.ttl = &ttl
		cache.wttl = true
	} else {
		ci.ttl = nil
	}
	cache.data[key] = ci
}

func (cache *memoryCache) Get(key string) (val any, ok bool) {
	cache.Lock()
	defer cache.Unlock()
	ci := cache.data[key]
	if ci == nil {
		return nil, false
	}
	ci.used = time.Now().UnixMilli()
	ci.hits++
	cache.data[key] = ci
	return ci.data, true
}
