package queue

import "sync"

// Node represents a single node in the queue
type Node struct {
	Val  interface{}
	Next *Node
	// _pad [5]int
}

// NewNode is a constructor for the Node type
func NewNode(val interface{}) *Node {
	n := new(Node)
	n.Val = val
	n.Next = nil

	return n
}

// LockFreeQueue
type LockFreeQueue struct {
	Head *Node
	Tail *Node
	Size uint32
	// _pad [5]int
}

// NewLockFreeQueue is a constructor for
// the LockFreeQueue type
func NewLockFreeQueue() *LockFreeQueue {
	queue := new(LockFreeQueue)
	queue.Head = new(Node)
	queue.Head.Next = nil
	queue.Tail = queue.Head
	queue.Size = 0

	return queue
}

// LockQueue is a thread safe
// queue that is not lock free
type LockQueue struct {
	Head *Node
	Tail *Node
	Size uint32
	lock sync.Mutex
}

// NewLockQueue is a constructor for the
// LockQueue type.
func NewLockQueue() *LockQueue {
	queue := new(LockQueue)
	queue.Head = new(Node)
	queue.Head.Next = nil
	queue.Tail = queue.Head
	queue.Size = 0

	return queue
}
