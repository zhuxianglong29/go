package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	// wg用来等待程序完成
	var wg sync.WaitGroup

	// 分配一个逻辑处理器给调度器使用
	runtime.GOMAXPROCS(2)

	// 计数加2，表示要等待两个goroutine
	wg.Add(2)

	// 创建两个goroutine
	fmt.Println("Create Goroutines")
	go printPrime("A", &wg)
	go printPrime("B", &wg)

	// 等待goroutine结束
	fmt.Println("Waiting To Finish")
	wg.Wait()

	fmt.Println("Terminating Program")
}

// printPrime 显示5000以内的素数值
func printPrime(prefix string, wg *sync.WaitGroup) {
	// 在函数退出时调用Done来通知main函数工作已经完成
	defer wg.Done()

next:
	for outer := 2; outer < 5000; outer++ {
		for inner := 2; inner < outer; inner++ {
			if outer%inner == 0 {
				continue next
			}
		}
		fmt.Printf("%s:%d\n", prefix, outer)
	}
	fmt.Println("Completed", prefix)
}
