package queue

// Node represents a single node in the queue
type Node struct {
	Val  interface{}
	Next *Node
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
