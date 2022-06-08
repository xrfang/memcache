package main

import (
	"fmt"
	"time"

	memorycache "github.com/larch0718/memory-cache"
)

func main() {
	cache := memorycache.New(nil)
	fmt.Println("设置值", "zhangsan")
	cache.Set("key1", "zhangsan", 2*time.Second)
	if it, ok := cache.Get("key1"); ok {
		fmt.Printf("#1 取到值 %v\n", it)
	} else {
		fmt.Println("获取值失败")
	}
	fmt.Println("sleep 5s")
	time.Sleep(5 * time.Second)
	if it, ok := cache.Get("key1"); ok {
		fmt.Printf("#2 取到值 %v\n", it)
	} else {
		fmt.Println("获取值失败")
	}
}
