package stack

import (
	"sync/atomic"
	"unsafe"
)

func (stack *LockFreeStack) Pop() interface{} {

	var oldTop *Item
	var newTop *Item

	for {
		oldTop = (*Item)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&stack.Top))))

		if oldTop == nil {
			continue
		}

		newTop = (*Item)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&oldTop.Next))))

		if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&stack.Top)), unsafe.Pointer(oldTop), unsafe.Pointer(newTop)) {
			break
		}
	}

	return oldTop.Value
}

func (stack *LockFreeStack) Push(value interface{}) {
	var oldTop *Item
	newTop := &Item{value, nil}

	for {
		oldTop = (*Item)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&stack.Top))))
		newTop.Next = oldTop

		if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&stack.Top)), unsafe.Pointer(oldTop), unsafe.Pointer(newTop)) {
			break
		}
	}

}

func (stack *LockFreeStack) Peek() interface{} {
	for {
		top := (*Item)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&stack.Top))))
		topValue := top.Value
		if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&stack.Top)), unsafe.Pointer(top), unsafe.Pointer(top)) {
			return topValue
		}
	}
}

func (stack *LockStack) Pop() interface{} {
	stack.accessLock.Lock()
	defer stack.accessLock.Unlock()

	oldTop := stack.Top
	newTop := oldTop.Next
	stack.Top = newTop

	return oldTop.Value
}

func (stack *LockStack) Push(value interface{}) {
	stack.accessLock.Lock()
	defer stack.accessLock.Unlock()

	oldTop := stack.Top
	newTop := &Item{value, nil}

	newTop.Next = oldTop
	stack.Top = newTop
}

func (stack *LockStack) Peek() interface{} {
	stack.accessLock.RLock()
	defer stack.accessLock.RUnlock()

	top := stack.Top

	if top == nil {
		return nil
	}
	return top.Value

}
