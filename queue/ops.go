package queue

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

// TODO: profile
// TODO: interleave inserts and deltes in all ds.

// Enqueue enqueues a `val`
func (q *LockFreeQueue) Enqueue(val interface{}) {
	newNode := NewNode(val)
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
	var newHead *Node
	for {
		pointer = (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&q.Head))))
		if atomic.CompareAndSwapPointer(
			(*unsafe.Pointer)(unsafe.Pointer(&pointer.Next)),
			unsafe.Pointer(nil),
			unsafe.Pointer(nil),
		) {
			return nil, fmt.Errorf("empty queue")
		}
		newHead = (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&pointer.Next))))
		ok := atomic.CompareAndSwapPointer(
			(*unsafe.Pointer)(unsafe.Pointer(&q.Head)),
			unsafe.Pointer(pointer),
			unsafe.Pointer(newHead),
		)
		if ok {
			break
		}
	}
	// decrement queue size by 1.
	atomic.AddUint32(&q.Size, ^uint32(0))
	return newHead.Val, nil
}

// Enqueue implements the LockQueue type, enqueues
// an element into a LockQueue. Makes use of mutex
func (q *LockQueue) Enqueue(val interface{}) {
	q.lock.Lock()
	defer q.lock.Unlock()

	newNode := NewNode(val)
	q.Tail.Next = newNode
	q.Tail = newNode
	q.Size++
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
	q.Size--
	if pointer.Next == nil {
		return nil, fmt.Errorf("empty queue")
	}

	return pointer.Next.Val, nil
}
