package main

import (
	"fmt"
	"sync"
	"time"
)

var numb int = 0
var mutex sync.Mutex

// var rwMutex sync.RWMutex //定义一个读写互斥锁
func add() {
	for i := 0; i < 10000; i++ {
		mutex.Lock()
		//mutex.Lock()
		numb = numb + 1
		mutex.Unlock()
		//mutex.Unlock()
		fmt.Println("numb=", numb)
		time.Sleep(time.Millisecond)
	}

}
func main() {

	go add()

	go add()

	for {

	}

}
