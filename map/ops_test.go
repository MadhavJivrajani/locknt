package lfmap

import (
	"math/rand"
	"sync"
	"testing"
)

func BenchmarkLockFreeDelAndIns(b *testing.B) {
	rand.Seed(420)
	m := NewMap(int64(b.N) + 1)
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
