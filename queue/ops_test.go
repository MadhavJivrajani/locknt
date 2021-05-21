package queue

import (
	"math/rand"
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
	for i := 0; i < b.N+1; i++ {
		q.Enqueue(i)
	}
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
	for i := 0; i < b.N+1; i++ {
		q.Enqueue(i)
	}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			q.Dequeue()
			wg.Done()
		}(&wg)
	}
	wg.Wait()
}

func BenchmarkLockFreeDelAndIns(b *testing.B) {
	rand.Seed(420)
	q := NewLockFreeQueue()
	var wg sync.WaitGroup
	for i := 0; i < b.N/2; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			q.Enqueue(int64(i))
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()
	for i := b.N / 2; i < b.N; i++ {
		wg.Add(2)
		go func(i int, wg *sync.WaitGroup) {
			q.Enqueue(int64(i))
			wg.Done()
		}(i, &wg)
		go func(i int, wg *sync.WaitGroup) {
			q.Dequeue()
			wg.Done()
		}(rand.Intn(int(i+1)), &wg)
	}
	wg.Wait()
}
