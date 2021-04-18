package stack

import (
	"sync/atomic"
	"unsafe"
)

func (stack *LockFreeStack) Pop() ValueType {

	var oldTop *Item
	var newTop *Item

	for {
		oldTop = (*Item)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&stack.Top))))

		if oldTop == nil {
			continue
		}

		newTop = (*Item)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&stack.Top.Next))))

		if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&stack.Top)), unsafe.Pointer(oldTop), unsafe.Pointer(newTop)) {
			break
		}
	}

	return oldTop.Value
}

func (stack *LockFreeStack) Push(value ValueType) {
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

func (stack *LockFreeStack) Peek(value ValueType) ValueType {
	for {
		top := (*Item)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&stack.Top))))
		topValue := top.Value
		if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&stack.Top)), unsafe.Pointer(top), unsafe.Pointer(top)) {
			return topValue
		}
	}
}

func (stack *LockStack) Pop() ValueType {

	stack.accessLock.Lock()

	oldTop := stack.Top
	newTop := oldTop.Next
	stack.Top = newTop

	stack.accessLock.Unlock()
	return oldTop.Value
}

func (stack *LockStack) Push(value ValueType) {

	stack.accessLock.Lock()

	oldTop := stack.Top
	newTop := &Item{value, nil}

	newTop.Next = oldTop
	stack.Top = newTop

	stack.accessLock.Unlock()
}

func (stack *LockStack) Peek(value ValueType) ValueType {
	stack.accessLock.RLock()

	top := stack.Top

	stack.accessLock.RUnlock()

	if top == nil {
		return nil
	}
	return top.Value

}
