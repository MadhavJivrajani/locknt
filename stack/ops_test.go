package stack

import (
	"sync"
	"testing"
)

func BenchmarkLockFreePush(b *testing.B) {
	s := NewLockFreeStack()
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			s.Push(i)
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()
}

func BenchmarkLockPush(b *testing.B) {
	s := NewLockStack()
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			s.Push(i)
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()
}

func BenchmarkLockFreePop(b *testing.B) {
	s := NewLockFreeStack()
	var wg sync.WaitGroup
	for i := 0; i < b.N+1; i++ {
		s.Push(i)
	}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			s.Pop()
			wg.Done()
		}(&wg)
	}
	wg.Wait()
}

func BenchmarkLockDequeue(b *testing.B) {
	s := NewLockStack()
	var wg sync.WaitGroup
	for i := 0; i < b.N+1; i++ {
		s.Push(i)
	}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			s.Pop()
			wg.Done()
		}(&wg)
	}
	wg.Wait()
}
