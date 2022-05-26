package main

import (
	"fmt"
	"time"

	memorycache "github.com/larch0718/memory-cache"
)

func main() {
	cache := memorycache.InitCache(memorycache.CacheOption{MaxItems: 100, Expire: 1 * time.Hour})
	fmt.Println("设置值", "zhangsan")
	cache.Add("key1", "zhangsan", 5*time.Second)
	if it, ok := cache.Get("key1"); ok {
		fmt.Printf("#1 取到值 %v\n", it)
	} else {
		fmt.Println("获取值失败")
	}
	fmt.Println("sleep 10s")
	time.Sleep(10 * time.Second)
	if it, ok := cache.Get("key1"); ok {
		fmt.Printf("#2 取到值 %v\n", it)
	} else {
		fmt.Println("获取值失败")
	}
}
