package list

import (
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"
)

type Node struct {
	Next      *Node
	Data      int64
	isDeleted int64
	_pad      [5]int
}

// Tis an ordered list
type LockFreeList struct {
	Head *Node
	Size uint32
	_pad [6]int
}

type LockList struct {
	Head *Node
	Size uint32
	lock sync.RWMutex
}

func newNode(data int64) *Node {
	new_node := &Node{
		Data:      data,
		Next:      nil,
		isDeleted: 0,
	}
	return new_node
}

func NewLockFreeList() *LockFreeList {
	return &LockFreeList{}
}

func NewLockList() *LockList {
	return &LockList{}
}

func ListCmp(ele1, ele2 int64) int64 {
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
		if ListCmp(headNode.Data, node.Data) > 0 {
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
			// headNext := (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&headNode.Next))))
			// if atomic.CompareAndSwapInt64(((*int64)(&headNext.isDeleted)), int64(1), int64(1)) {
			// 	continue
			// }
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
		if ListCmp(nextNode.Data, node.Data) > 0 {
			headNext := (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&headNode.Next))))
			if atomic.CompareAndSwapInt64(((*int64)(&headNext.isDeleted)), int64(1), int64(1)) {
				continue
			}
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

		if ListCmp(nextNode.Data, node.Data) == 0 {
			return fmt.Errorf("Element exists")
		}

		headNode = nextNode
	}
}

func (list *LockList) Insert(data int64) error {
	list.lock.Lock()
	defer list.lock.Unlock()
	node := newNode(data)
	head := list.Head

	if head == nil {
		list.Head = node
		list.Size++
		return nil
	}
	headNode := list.Head
	if ListCmp(headNode.Data, node.Data) > 0 {
		node.Next = headNode
		list.Head = node
		list.Size++
		return nil
	}

	nextNode := headNode.Next

	for nextNode != nil {
		if ListCmp(nextNode.Data, node.Data) > 0 {
			node.Next = nextNode
			headNode.Next = node
			list.Size++
			return nil
		}
		if ListCmp(nextNode.Data, node.Data) == 0 {
			return fmt.Errorf("Value already in list")
		}
		headNode = nextNode
		nextNode = nextNode.Next
	}

	headNode.Next = node
	node.Next = nextNode
	list.Size++
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

		if ListCmp(headNode.Data, data) == 0 {
			nextPtr := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&headNode.Next)))
			atomic.AddInt64(&headNode.isDeleted, 1)
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

		if ListCmp(nextNode.Data, data) > 0 {
			return fmt.Errorf("Value not found")
		}

		if ListCmp(nextNode.Data, data) == 0 {
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

func (list *LockList) Delete(data int64) error {

	list.lock.Lock()
	defer list.lock.Unlock()
	head := list.Head
	if head == nil {
		return fmt.Errorf("List empty")
	}

	headNode := list.Head

	if ListCmp(headNode.Data, data) == 0 {
		list.Head = list.Head.Next
		list.Size--
		return nil
	}

	nextPtr := headNode.Next

	if nextPtr == nil {
		return fmt.Errorf("Value not found")
	}

	nextNode := headNode.Next

	if ListCmp(nextNode.Data, data) > 0 {
		return fmt.Errorf("Value not found")
	}

	for nextNode != nil {
		if ListCmp(nextNode.Data, data) == 0 {
			replacementPtr := nextNode.Next
			headNode.Next = replacementPtr
			list.Size--
			return nil
		}
		headNode = nextNode
		nextNode = nextNode.Next
	}
	return nil
}

func PrintList(list *LockFreeList) {
	currNode := list.Head
	for currNode != nil {
		fmt.Println(currNode.Data)
		currNode = currNode.Next
	}
}
