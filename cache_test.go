package memorycache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	cache := InitCache(CacheOption{MaxItems: 100, Expire: 1 * time.Hour})
	cache.Add("key1", "zhangsan", 5*time.Second)
	if it, ok := cache.Get("key1"); ok {
		if it.(string) != "zhangsan" {
			t.Fatalf("执行错误，期望取到值: %v  实际取到值: %v\n", "zhangsan", it)
		}
	} else {
		t.Fatalf("执行错误，期望取到值: %v  实际未取到值\n", "zhangsan")
	}
	t.Log("第一阶段测试通过")
	time.Sleep(10 * time.Second)
	if it, ok := cache.Get("key1"); ok {
		t.Fatalf("执行错误，期望未取到值  实际未取到值: %v\n", it)
	}
	t.Log("第二阶段测试通过")
}
