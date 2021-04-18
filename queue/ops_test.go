package queue

import (
	"sync"
	"testing"
)

func BenchmarkLockFreeEnqueue(b *testing.B) {
	q := NewLockFreeQueue()
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			q.Enqueue(i)
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()
}

func BenchmarkLockEnqueue(b *testing.B) {
	q := NewLockQueue()
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			q.Enqueue(i)
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()
}

func BenchmarkLockFreeDequeue(b *testing.B) {
	q := NewLockFreeQueue()
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			q.Dequeue()
			wg.Done()
		}(&wg)
	}
	wg.Wait()
}

func BenchmarkLockDequeue(b *testing.B) {
	q := NewLockQueue()
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			q.Dequeue()
			wg.Done()
		}(&wg)
	}
	wg.Wait()
}
