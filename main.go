package main

import (
	"math/rand"
	"sync"

	"github.com/MadhavJivrajani/locknt/list"
)

func main() {

	// for i := 0; i < 10; i++ {
	// 	wg.Add(1)
	// 	go func(i int, wg *sync.WaitGroup) {
	// 		err := l.Insert(int64(i))
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		wg.Done()
	// 	}(i, &wg)
	// }
	// wg.Wait()
	// fmt.Println(l.Size)
	// for i := 0; i < 10; i++ {
	// 	wg.Add(1)
	// 	go func(i int, wg *sync.WaitGroup) {
	// 		err := l.Delete(int64(i))
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		wg.Done()
	// 	}(i, &wg)
	// }
	// wg.Wait()

	b := 20000
	rand.Seed(420)
	s := list.NewLockFreeList()
	var wg sync.WaitGroup
	for i := 0; i < b/2; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			s.Insert(int64(i))
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()
	for i := b / 2; i < b; i++ {
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
