package main

import (
	"sync"

	"github.com/MadhavJivrajani/locknt/list"
)

func main() {
	l := list.NewLockFreeList()
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			err := l.Insert(int64(i))
			if err != nil {
				panic(err)
			}
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()
}
