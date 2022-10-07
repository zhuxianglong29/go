package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"golang.org/x/sync/semaphore"
)

// Example_workerPool demonstrates how to use a semaphore to limit the number of
// goroutines working on parallel tasks.
//
// This use of a semaphore mimics a typical “worker pool” pattern, but without
// the need to explicitly shut down idle workers when the work is done.
func main() {
	ctx := context.TODO()

	var (
		sem = semaphore.NewWeighted(1)
	)
	var deadLockCase time.Duration = 1
	go func() {
		time.Sleep(time.Second * deadLockCase)
		fmt.Println("Writer Try")
		if err := sem.Acquire(ctx, 1); err != nil { //获取锁、如果是Acquire会等待，而TryAcquire不会等待
			log.Printf("Failed to acquire semaphore: %v", err)

		}
		fmt.Println("Writer Fetch")
		time.Sleep(time.Second * 1)
		fmt.Println("Writer Release")
		sem.Release(1)
	}()

	fmt.Println("Reader 1 Try")
	if err := sem.Acquire(ctx, 1); err != nil { //获取锁、如果是Acquire会等待，而TryAcquire不会等待
		log.Printf("Failed to acquire semaphore: %v", err)

	}
	fmt.Println("Reader 1 Fetch")

	time.Sleep(time.Second * 2)
	fmt.Println("Reader 2 Try")
	if err := sem.Acquire(ctx, 1); err != nil { //获取锁、如果是Acquire会等待，而TryAcquire不会等待
		log.Printf("Failed to acquire semaphore: %v", err)

	}
	fmt.Println("Reader 2 Fetch")

	time.Sleep(time.Second * 2)
	fmt.Println("Reader 1 Release")
	sem.Release(1)

	time.Sleep(time.Second * 1)
	fmt.Println("Reader 2 Release")
	sem.Release(1)
	time.Sleep(time.Second * 2)
	fmt.Println("Done")

}
