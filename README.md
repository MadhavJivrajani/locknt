# locknt
A collections of Lock-free data structures in Golang

## Queue
The queue implementation is based on [this](http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.53.8674&rep=rep1&type=pdf) paper.

### Ex
```go
package main

import (
	"sync"

	"github.com/MadhavJivrajani/locknt/queue"
)

func main() {
	q := queue.NewLockFreeQueue()
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			q.Enqueue(i)
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()
	queue.PrintQueue(q)
}
```
### Benchmarking and generating call graphs

```
go test -bench=^<benchmark_name>$ -benchtime=100000x -cpuprofile profile_file.out
go tool pprof profile_file.out
(pprof) web # insde pprof, type this
```
