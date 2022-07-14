package memcache

import (
	"testing"
	"time"
)

func TestGetSet(t *testing.T) {
	t.Log("测试设置与更新")
	cache := New(nil)
	cache.Set("name", "zhangsan")
	v, ok := cache.Get("name")
	if !ok || v.(string) != "zhangsan" {
		t.Fatalf("expect 'zhangsan', got %v", v)
	}
	cache.Set("name", "lisi")
	v, ok = cache.Get("name")
	if !ok || v.(string) != "lisi" {
		t.Fatalf("expect 'lisi', got %v", v)
	}
	t.Log("GET/SET测试通过")
}

func TestExpire(t *testing.T) {
	t.Log("测试超时机制")
	cache := New(nil)
	cache.Set("key1", "zhangsan", 2*time.Second)
	if it, ok := cache.Get("key1"); ok {
		if it.(string) != "zhangsan" {
			t.Fatalf("执行错误，期望取到值: %v  实际取到值: %v\n", "zhangsan", it)
		}
	} else {
		t.Fatalf("执行错误，期望取到值: %v  实际未取到值\n", "zhangsan")
	}
	t.Log("第一阶段测试通过")
	time.Sleep(4 * time.Second)
	if it, ok := cache.Get("key1"); ok {
		t.Fatalf("执行错误，期望未取到值  实际取到值: %v\n", it)
	}
	t.Log("第二阶段测试通过")
}

func TestLRU(t *testing.T) {
	t.Log("测试LRU清除策略")
	cache := New(&Option{Capacity: 2})
	cache.Set("key1", "value1")
	time.Sleep(time.Second)
	cache.Set("key2", "value2")
	_, ok := cache.Get("key1")
	if !ok {
		t.Fatal("key1不见了")
	}
	cache.Set("key3", "value3")
	time.Sleep(2 * time.Second)
	_, ok = cache.Get("key2")
	if ok {
		t.Fatal("key2还在")
	}
	_, ok1 := cache.Get("key1")
	_, ok3 := cache.Get("key3")
	if !ok1 || !ok3 {
		t.Fatalf("key1或key3不见了")
	}
	t.Log("LRU测试通过")
}

func TestLFU(t *testing.T) {
	t.Log("测试LFU清除策略")
	cache := New(&Option{Capacity: 2, Policy: PolicyLFU})
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	_, ok := cache.Get("key1")
	if !ok {
		t.Fatal("key1不见了")
	}
	cache.Set("key3", "value3")
	time.Sleep(2 * time.Second)
	_, ok = cache.Get("key2")
	if ok {
		t.Fatal("key2还在")
	}
	_, ok1 := cache.Get("key1")
	_, ok3 := cache.Get("key3")
	if !ok1 || !ok3 {
		t.Fatalf("key1或key3不见了")
	}
	t.Log("LFU测试通过")
}
