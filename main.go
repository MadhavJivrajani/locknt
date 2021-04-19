package main

import (
	"fmt"
	"sync"

	"github.com/MadhavJivrajani/locknt/list"
)

func main() {
	l := list.NewLockList()
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
	fmt.Println(l.Size)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			err := l.Delete(int64(i))
			if err != nil {
				panic(err)
			}
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()
}
