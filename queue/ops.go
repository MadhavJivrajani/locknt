package queue

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

// TODO: mem profile
// TODO: interleave inserts and deltes in all ds.

// Enqueue enqueues a `val`
func (q *LockFreeQueue) Enqueue(val interface{}) {
	newNode := &Node{val, nil}
	var pointer *Node
	for {
		pointer = (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&q.Tail))))
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
		unsafe.Pointer(pointer),
		unsafe.Pointer(newNode),
	)
	// increment queue size by 1.
	atomic.AddUint32(&q.Size, 1)
}

// Dequeue deques an element from the queue and returns it
// along with an error, if any. Errors arise when the queue
// is empty and dequeue is attempted.
func (q *LockFreeQueue) Dequeue() (interface{}, error) {
	var pointer *Node
	for {
		pointer = (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&q.Head))))
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
	q.lock.Lock()
	defer q.lock.Unlock()

	newNode := &Node{val, nil}
	q.Tail.Next = newNode
	q.Tail = newNode
	q.Size += 1
}

// Dequeue implements the LockQueue type, dequeues
// an element, returns the element if no error
// occured. Makes use of mutex
func (q *LockQueue) Dequeue() (interface{}, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	pointer := q.Head
	if pointer.Next == nil {
		return nil, fmt.Errorf("empty queue")
	}
	pointer = pointer.Next
	q.Size -= 1
	return pointer.Next.Val, nil
}
