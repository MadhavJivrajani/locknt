package LocklessStack

import (
	sync "sync"
)

type ValueType []interface{}

type Item struct {
	Value ValueType
	Next *Item
}


type Stack struct {
	Top *Item
	accessLock sync.RWMutex
}

func makeStack() *Stack {
	newStack := new(Stack)
	newStack.Top = nil
	return newStack
}

func Pop( stack *Stack ) ValueType {

	stack.accessLock.Lock();

	oldTop := stack.Top
	newTop := oldTop.Next
	stack.Top = newTop

	stack.accessLock.Unlock();
	return oldTop.Value
}

func Push( stack *Stack, value ValueType) {

	stack.accessLock.Lock();

	oldTop := stack.Top
	newTop := &Item{value, nil}

	newTop.Next = oldTop
	stack.Top = newTop
	
	stack.accessLock.Unlock();
}

func Peek (stack *Stack, value ValueType) ValueType {
	stack.accessLock.RLock();

	top := stack.Top

	stack.accessLock.RUnlock();

	if top == nil {
		return nil;		
	}
	return top.Value

}