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
	Size uint32
}

type LockedList struct {
	Head *Node
	Tail *Node
	Size uint32
	lock sync.RWMutex
}

func newNode(data int64) *Node {
	new_node := &Node{
		Data: data,
		Next: nil,
	}
	return new_node
}

func NewLockFreeList() *LockFreeList {
	return &LockFreeList{}
}

func NewLockedList() *LockedList {
	list := new(LockedList)
	list.Head = new(Node)
	list.Tail = new(Node)
	list.Head.Next = list.Tail
	list.Tail.Next = nil
	list.Size = 0

	return list
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

func LockFreeListCmp(ele1, ele2 int64) int64 {
	return ele1 - ele2
}

func (list *LockFreeList) Insert(data int64) error {

	node := newNode(data)

	var headPtr unsafe.Pointer
	var headNode *Node

	for {
		headPtr = atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&list.Head)))

		if headPtr == nil {
			ok := atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&list.Head)),
				headPtr,
				unsafe.Pointer(node))
			if !ok {
				continue
			}

			atomic.AddUint32(&list.Size, 1)
			return nil
		}

		headNode = (*Node)(headPtr)
		if LockFreeListCmp(headNode.Data, node.Data) > 0 {
			node.Next = (*Node)(headPtr)
			ok := atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&list.Head)),
				headPtr,
				unsafe.Pointer(node))
			if !ok {
				continue
			}

			atomic.AddUint32(&list.Size, 1)
			return nil
		}
		break
	}

	for {
		nextPtr := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&headNode.Next)))
		if nextPtr == nil {
			ok := atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&headNode.Next)),
				nextPtr,
				unsafe.Pointer(node))
			if !ok {
				continue
			}

			atomic.AddUint32(&list.Size, 1)
			return nil
		}

		nextNode := (*Node)(nextPtr)
		if LockFreeListCmp(nextNode.Data, node.Data) > 0 {
			node.Next = (*Node)(nextPtr)
			ok := atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&headNode.Next)),
				nextPtr,
				unsafe.Pointer(node))
			if !ok {
				continue
			}

			atomic.AddUint32(&list.Size, 1)
			return nil
		}

		if LockFreeListCmp(nextNode.Data, node.Data) == 0 {
			return fmt.Errorf("Element exists")
		}

		headNode = nextNode
	}
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
	var headPtr unsafe.Pointer
	var headNode *Node
	for {
		headPtr = atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&list.Head)))

		if headPtr == nil {
			return fmt.Errorf("Value not found")
		}

		headNode = (*Node)(headPtr)

		if LockFreeListCmp(headNode.Data, data) == 0 {
			nextPtr := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&headNode.Next)))
			ok := atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&list.Head)),
				headPtr,
				nextPtr)
			if !ok {
				continue
			}

			atomic.AddUint32(&list.Size, ^uint32(0))
			return nil
		}
		break
	}

	for {
		nextPtr := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&headNode.Next)))
		if nextPtr == nil {
			return fmt.Errorf("Value not found")
		}

		nextNode := (*Node)(nextPtr)

		if LockFreeListCmp(nextNode.Data, data) > 0 {
			return fmt.Errorf("Value not found")
		}

		if LockFreeListCmp(nextNode.Data, data) == 0 {
			replacementPtr := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&nextNode.Next)))
			ok := atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&headNode.Next)),
				nextPtr,
				replacementPtr)
			if !ok {
				continue
			}

			atomic.AddUint32(&list.Size, ^uint32(0))
			return nil
		}

		headNode = nextNode
	}
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
