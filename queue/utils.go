package queue

import "fmt"

func PrintQueue(q *LockFreeQueue) {
	p := q.Head.Next
	for p != nil {
		fmt.Println(p.Val)
		p = p.Next
	}
}
