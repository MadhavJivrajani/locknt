package queue

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

// Enqueue enqueues a `val`
func (q *LockFreeQueue) Enqueue(val interface{}) {
	// create the node to be inserted.
	newNode := NewNode(val)
	pointer := q.Tail
	oldPointer := pointer

	for {
		for pointer.Next != nil {
			pointer = pointer.Next
		}
		ok := atomic.CompareAndSwapPointer(
			(*unsafe.Pointer)(unsafe.Pointer(&pointer.Next)),
			unsafe.Pointer(nil),
			unsafe.Pointer(newNode),
		)
		if ok {
			break
		}
	}
	atomic.CompareAndSwapPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&q.Tail)),
		unsafe.Pointer(oldPointer),
		unsafe.Pointer(newNode),
	)
	// increment queue size.
	atomic.AddUint32(&q.Size, 1)
}

// Dequeue deques an element from the queue and returns it
// along with an error, if any. Errors arise when the queue
// is empty and dequeue is attempted.
func (q *LockFreeQueue) Dequeue() (interface{}, error) {
	var pointer *Node
	for {
		pointer = q.Head
		if pointer.Next == nil {
			return nil, fmt.Errorf("empty queue")
		}
		ok := atomic.CompareAndSwapPointer(
			(*unsafe.Pointer)(unsafe.Pointer(&q.Head)),
			unsafe.Pointer(pointer),
			unsafe.Pointer(pointer.Next),
		)
		if ok {
			break
		}
	}
	// decrement queue size by 1.
	atomic.AddUint32(&q.Size, ^uint32(0))
	return pointer.Next.Val, nil
}

// Enqueue implements the LockQueue type, enqueues
// an element into a LockQueue. Makes use of mutex
func (q *LockQueue) Enqueue(val interface{}) {

}
