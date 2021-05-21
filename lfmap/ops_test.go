package lfmap

import (
	"math/rand"
	"sync"
	"testing"
)

func BenchmarkLockFreeLookupAndIns(b *testing.B) {
	rand.Seed(420)
	m := NewLockFreeMap(int64(b.N) + 1)
	var wg sync.WaitGroup
	for i := 0; i < b.N/2; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			m.Insert(int64(i), rand.Intn(int(i+1)))
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()
	for i := b.N / 2; i < b.N; i++ {
		wg.Add(2)
		go func(i int, wg *sync.WaitGroup) {
			m.Insert(int64(i), rand.Intn(int(i+1)))
			wg.Done()
		}(i, &wg)
		go func(i int, wg *sync.WaitGroup) {
			m.Lookup(int64(i))
			wg.Done()
		}(rand.Intn(int(i+1)), &wg)
	}
	wg.Wait()
}

func BenchmarkLockLookupAndIns(b *testing.B) {
	rand.Seed(420)
	m := NewLockMap(int64(b.N) + 1)
	var wg sync.WaitGroup
	for i := 0; i < b.N/2; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			m.Insert(int64(i), rand.Intn(int(i+1)))
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()
	for i := b.N / 2; i < b.N; i++ {
		wg.Add(2)
		go func(i int, wg *sync.WaitGroup) {
			m.Insert(int64(i), rand.Intn(int(i+1)))
			wg.Done()
		}(i, &wg)
		go func(i int, wg *sync.WaitGroup) {
			m.Lookup(int64(i))
			wg.Done()
		}(rand.Intn(int(i+1)), &wg)
	}
	wg.Wait()
}

func BenchmarkLockFreeLookup(b *testing.B) {
	rand.Seed(420)
	m := NewLockFreeMap(int64(b.N) + 1)
	var wg sync.WaitGroup
	for i := 0; i < b.N/2; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			m.Insert(int64(i), rand.Intn(int(i+1)))
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()
	for i := b.N / 2; i < b.N; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			m.Lookup(int64(i))
			wg.Done()
		}(rand.Intn(int(i+1)), &wg)
	}
	wg.Wait()
}

func BenchmarkLockLookup(b *testing.B) {
	rand.Seed(420)
	m := NewLockMap(int64(b.N) + 1)
	var wg sync.WaitGroup
	for i := 0; i < b.N/2; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			m.Insert(int64(i), rand.Intn(int(i+1)))
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()
	for i := b.N / 2; i < b.N; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			m.Lookup(int64(i))
			wg.Done()
		}(rand.Intn(int(i+1)), &wg)
	}
	wg.Wait()
}

func BenchmarkLockFreeIns(b *testing.B) {
	rand.Seed(420)
	m := NewLockFreeMap(int64(b.N) + 1)
	var wg sync.WaitGroup
	for i := 0; i < b.N/2; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			m.Insert(int64(i), rand.Intn(int(i+1)))
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()
	for i := b.N / 2; i < b.N; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			m.Insert(int64(i), rand.Intn(int(i+1)))
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()
}

func BenchmarkLockIns(b *testing.B) {
	rand.Seed(420)
	m := NewLockMap(int64(b.N) + 1)
	var wg sync.WaitGroup
	for i := 0; i < b.N/2; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			m.Insert(int64(i), rand.Intn(int(i+1)))
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()
	for i := b.N / 2; i < b.N; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			m.Insert(int64(i), rand.Intn(int(i+1)))
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()
}

func BenchmarkLockFreeInsertIfDoesntExistAndLookup(b *testing.B) {
	rand.Seed(420)
	m := NewLockFreeMap(int64(b.N) + 1)
	var wg sync.WaitGroup
	for i := 0; i < b.N/2; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			m.InsertIfDoesntExist(int64(i), rand.Intn(int(i+1)))
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()
	for i := b.N / 2; i < b.N; i++ {
		wg.Add(2)
		go func(i int, wg *sync.WaitGroup) {
			m.InsertIfDoesntExist(int64(i), rand.Intn(int(i+1)))
			wg.Done()
		}(i, &wg)
		go func(i int, wg *sync.WaitGroup) {
			m.Lookup(int64(i))
			wg.Done()
		}(rand.Intn(int(i+1)), &wg)
	}
	wg.Wait()
}

func cmp(a interface{}, b interface{}) bool {
	return (a.(int) % 100) > (b.(int) % 100)
}
func BenchmarkLockFreeInsertCompareAndLookup(b *testing.B) {

	rand.Seed(420)
	m := NewLockFreeMap(int64(b.N) + 1)
	var wg sync.WaitGroup
	for i := 0; i < b.N/2; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			m.InsertCompare(int64(i), rand.Intn(int(i+1)), cmp)
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()
	for i := b.N / 2; i < b.N; i++ {
		wg.Add(2)
		go func(i int, wg *sync.WaitGroup) {
			m.InsertCompare(int64(i), rand.Intn(int(i+1)), cmp)
			wg.Done()
		}(i, &wg)
		go func(i int, wg *sync.WaitGroup) {
			m.Lookup(int64(i))
			wg.Done()
		}(rand.Intn(int(i+1)), &wg)
	}
	wg.Wait()
}
