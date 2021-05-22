package list

import (
	"math/rand"
	"runtime"
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
		}(rand.Intn(i+1), &wg)
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
		}(rand.Intn(i+1), &wg)
	}
	wg.Wait()
}

func generateRandom(size int) []int {
	var tests []int
	rand.Seed(420)
	for i := 0; i < size; i++ {
		tests = append(tests, i+1)
	}
	r := rand.New(rand.NewSource(420))
	r.Shuffle(len(tests), func(i, j int) { tests[i], tests[j] = tests[j], tests[i] })

	return tests
}

func BenchmarkLockFreeDelAndIns(b *testing.B) {
	runtime.GOMAXPROCS(1)
	rand.Seed(420)
	s := NewLockFreeList()
	var wg sync.WaitGroup
	for i := 0; i < b.N/2; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			s.Insert(int64(i))
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()
	for i := b.N / 2; i < b.N; i++ {
		wg.Add(2)
		go func(i int, wg *sync.WaitGroup) {
			s.Insert(int64(i))
			wg.Done()
		}(i, &wg)
		go func(i int, wg *sync.WaitGroup) {
			s.Delete(int64(i))
			wg.Done()
		}(rand.Intn(int(i+1)), &wg)
	}
	wg.Wait()
}

func BenchmarkLockDelAndIns(b *testing.B) {
	runtime.GOMAXPROCS(1)
	rand.Seed(420)
	s := NewLockList()
	var wg sync.WaitGroup
	for i := 0; i < b.N/2; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			s.Insert(int64(i))
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()
	for i := b.N / 2; i < b.N; i++ {
		wg.Add(2)
		go func(i int, wg *sync.WaitGroup) {
			s.Insert(int64(i))
			wg.Done()
		}(i, &wg)
		go func(i int, wg *sync.WaitGroup) {
			s.Delete(int64(i))
			wg.Done()
		}(rand.Intn(int(i+1)), &wg)
	}
	wg.Wait()
}
