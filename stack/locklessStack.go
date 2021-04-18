package LocklessStack

import (
	"sync/atomic"
	"unsafe"
)

type ValueType []interface{}

type Item struct {
	Value ValueType
	Next *Item
}


type Stack struct {
	Top *Item
}

func makeStack() *Stack {
	newStack := new(Stack)
	newStack.Top = nil
	return newStack
}

func Pop( stack *Stack ) ValueType {

	var oldTop *Item
	var newTop *Item

	for {
		oldTop = (*Item) (atomic.LoadPointer( (*unsafe.Pointer )(unsafe.Pointer(&stack.Top))))

		if oldTop == nil {
			continue
		}

		newTop = (*Item) (atomic.LoadPointer( (*unsafe.Pointer )(unsafe.Pointer(&stack.Top.Next))))

		if atomic.CompareAndSwapPointer( (*unsafe.Pointer )(unsafe.Pointer(&stack.Top)), unsafe.Pointer(oldTop), unsafe.Pointer(newTop) ) {
			break
		}
	}

	return oldTop.Value
}

func Push( stack *Stack, value ValueType) {
	var oldTop *Item
	newTop := &Item{value, nil}

	for {
		oldTop = (*Item) (atomic.LoadPointer(( *unsafe.Pointer ) (unsafe.Pointer(&stack.Top)) ))
		newTop.Next = oldTop

		if atomic.CompareAndSwapPointer( (*unsafe.Pointer )(unsafe.Pointer(&stack.Top)), unsafe.Pointer(oldTop), unsafe.Pointer(newTop) ) {
			break
		}
	}

} 

func Peek (stack *Stack, value ValueType) ValueType {
	for {
		top := (*Item) (atomic.LoadPointer(( *unsafe.Pointer )(unsafe.Pointer(&stack.Top))))
		topValue := top.Value
		if atomic.CompareAndSwapPointer( (*unsafe.Pointer )(unsafe.Pointer(&stack.Top)), unsafe.Pointer(top), unsafe.Pointer(top) ) {
			return topValue;
		}
	}
}
