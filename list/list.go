package list

import (
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"
)

type Node struct {
	Next *Node
	Data int64
}

// Tis an ordered list
type LockFreeList struct {
	Head *Node
	Tail *Node
	Size uint32
}

type LockedList struct {
	Head *Node
	Tail *Node
	Size uint32
	lock sync.RWMutex
}

func newNode(data int64) *Node {
	new_node := new(Node)
	new_node.Data = data
	new_node.Next = nil
	return new_node
}

func NewLockFreeList() *LockFreeList {
	list := new(LockFreeList)
	list.Head = new(Node)
	list.Tail = new(Node)
	list.Head.Next = list.Tail
	list.Size = 0

	return list
}

func NewLockedList() *LockedList {
	list := new(LockedList)
	list.Head = new(Node)
	list.Tail = new(Node)
	list.Head.Next = list.Tail
	list.Size = 0

	return list
}

// Lockfree List is not using 2 stage deletion as of yet -> may lead to some nasty cases

// Required functions
// insert - have an arg that tells whether the list is ordered (based on key)
// delete - based on key
// find
// search - required for ordered list insert

func (list *LockFreeList) search(data int64, left_node **Node) *Node {
	var left_node_next *Node
	var right_node *Node

	for {
		t := list.Head
		t_next := t.Next

		for t.Data < data {
			*left_node = t
			left_node_next = t_next
			t = t_next
			if t == list.Tail {
				break
			}
			t_next = t.Next
		}
		right_node = t

		if left_node_next == right_node {
			return right_node
		}
	}
}

func (list *LockedList) search(data int64, left_node **Node) *Node {
	var left_node_next *Node
	var right_node *Node

	for {
		t := list.Head
		t_next := t.Next

		for t.Data < data {
			*left_node = t
			left_node_next = t_next
			t = t_next
			if t == list.Tail {
				break
			}
			t_next = t.Next
		}
		right_node = t

		if left_node_next == right_node {
			return right_node
		}
	}
}

func (list *LockFreeList) Insert(data int64) error {

	node := newNode(data)
	var right_node *Node
	var left_node *Node

	for {
		right_node = list.search(data, &left_node)
		if right_node != list.Tail && right_node.Data == data {
			return fmt.Errorf("Key already exists in list")
		}
		node.Next = right_node

		ok := atomic.CompareAndSwapPointer(
			(*unsafe.Pointer)(unsafe.Pointer(left_node.Next)),
			unsafe.Pointer(right_node),
			unsafe.Pointer(node))

		if ok {
			break
		}
	}
	atomic.AddUint32(&list.Size, 1)
	return nil
}

func (list *LockedList) Insert(data int64) error {
	node := newNode(data)
	list.lock.Lock()

	var left_node *Node
	right_node := list.search(data, &left_node)
	if right_node != list.Tail && right_node.Data == data {
		return fmt.Errorf("Key already exists in list")
	}
	node.Next = right_node
	left_node.Next = node

	list.Size++

	list.lock.Unlock()
	return nil
}

func (list *LockFreeList) Delete(data int64) error {
	var right_node *Node
	var right_node_next *Node
	var left_node *Node

	for {
		right_node = list.search(data, &left_node)
		if right_node == list.Tail || right_node.Data != data {
			return fmt.Errorf("Data not present in list")
		}
		right_node_next = right_node.Next
		ok := atomic.CompareAndSwapPointer(
			(*unsafe.Pointer)(unsafe.Pointer(&left_node.Next)),
			(unsafe.Pointer(right_node)),
			(unsafe.Pointer(right_node_next)))
		if ok {
			break
		}
	}
	atomic.AddUint32(&list.Size, ^uint32(0))
	return nil
}

func (list *LockedList) Delete(data int64) error {
	var left_node *Node

	list.lock.Lock()

	right_node := list.search(data, &left_node)
	if right_node == list.Tail || right_node.Data != data {
		return fmt.Errorf("Data not present in list")
	}
	left_node.Next = right_node.Next
	list.Size--

	list.lock.Unlock()
	return nil
}

func PrintList(list *LockFreeList) {
	currNode := list.Head
	for currNode != nil {
		fmt.Println(currNode.Data)
		currNode = currNode.Next
	}
}
