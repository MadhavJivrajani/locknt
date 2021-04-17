package main

import (
	"sync"

	"github.com/MadhavJivrajani/locknt/queue"
)

func main() {
	q := queue.NewLockFreeQueue()
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			q.Enqueue(i)
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()
	queue.PrintQueue(q)
}
