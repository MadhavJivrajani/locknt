package stack

import (
	"math/rand"
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

func BenchmarkLockPop(b *testing.B) {
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

func BenchmarkLockFreeDelAndIns(b *testing.B) {
	rand.Seed(420)
	s := NewLockFreeStack()
	var wg sync.WaitGroup
	for i := 0; i < b.N/2; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			s.Push(int64(i))
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()
	for i := b.N / 2; i < b.N; i++ {
		wg.Add(2)
		go func(i int, wg *sync.WaitGroup) {
			s.Push(int64(i))
			wg.Done()
		}(i, &wg)
		go func(i int, wg *sync.WaitGroup) {
			s.Pop()
			wg.Done()
		}(rand.Intn(int(i+1)), &wg)
	}
	wg.Wait()
}

func BenchmarkLockDelAndIns(b *testing.B) {
	rand.Seed(420)
	s := NewLockStack()
	var wg sync.WaitGroup
	for i := 0; i < b.N/2; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			s.Push(int64(i))
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()
	for i := b.N / 2; i < b.N; i++ {
		wg.Add(2)
		go func(i int, wg *sync.WaitGroup) {
			s.Push(int64(i))
			wg.Done()
		}(i, &wg)
		go func(i int, wg *sync.WaitGroup) {
			s.Pop()
			wg.Done()
		}(rand.Intn(int(i+1)), &wg)
	}
	wg.Wait()
}
