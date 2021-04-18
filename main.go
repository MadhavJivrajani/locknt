package main

import (
	"sync"

	"github.com/MadhavJivrajani/locknt/queue"
)

func main() {
	q := queue.NewLockQueue()
	var wg sync.WaitGroup
	//s := stack.NewLockFreeStack()

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			q.Enqueue(i)
			//s.Push(stack.ValueType{i})
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			q.Dequeue()
			//s.Pop()
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()
	//queue.PrintQueue(q)
}
