package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	rw := new(sync.RWMutex)
	var deadLockCase time.Duration = 1
	go func() {
		time.Sleep(time.Second * deadLockCase)
		fmt.Println("Writer Try")
		rw.Lock()
		fmt.Println("Writer Fetch")
		time.Sleep(time.Second * 1)
		fmt.Println("Writer Release")
		rw.Unlock()
	}()

	fmt.Println("Reader 1 Try")
	rw.RLock()
	fmt.Println("Reader 1 Fetch")

	time.Sleep(time.Second * 2)
	fmt.Println("Reader 2 Try")
	rw.RLock()
	fmt.Println("Reader 2 Fetch")

	time.Sleep(time.Second * 2)
	fmt.Println("Reader 1 Release")
	rw.RUnlock()

	time.Sleep(time.Second * 1)
	fmt.Println("Reader 2 Release")
	rw.RUnlock()
	time.Sleep(time.Second * 2)
	fmt.Println("Done")
}
