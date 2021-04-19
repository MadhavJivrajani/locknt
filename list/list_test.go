package list

import (
	"sync"
	"testing"
)

func BenchmarkLockFreeInsert(b *testing.B) {
	s := NewLockFreeList()
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			s.Insert(int64(i))
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()
}

func BenchmarkLockInsert(b *testing.B) {
	s := NewLockList()
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			s.Insert(int64(i))
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()
}

func BenchmarkLockFreeDelete(b *testing.B) {
	s := NewLockFreeList()
	var wg sync.WaitGroup
	for i := 0; i < b.N+1; i++ {
		s.Insert(int64(i))
	}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			s.Delete(int64(i))
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()
}

func BenchmarkLockDelete(b *testing.B) {
	s := NewLockList()
	var wg sync.WaitGroup
	for i := 0; i < b.N+1; i++ {
		s.Insert(int64(i))
	}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			s.Delete(int64(i))
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()
}
