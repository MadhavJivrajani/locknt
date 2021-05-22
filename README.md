# locknt
A collection of lock-free data structures in Golang, done as a project under the course: Heterogenous Parallelism (UE18CS342), at PES University.

## Data structures implemented:
- [Queue](./queue)
- [Stack](./stack)
- [List](./list)
- [Map](./lfmap)

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
		go func(i int) {
			defer wg.Done()
			q.Enqueue(i)
			
		}(i)
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

### Work done by:
- Madhav Jivrajani
- M S Akshatha Laxmi
- Sparsh Temani

### Presentations
The presentations done for this project can be found [here](./assets).
