package stack

import sync "sync"

// Item
type Item struct {
	Value interface{}
	Next  *Item
}

// LockFreeStack
type LockFreeStack struct {
	Top *Item
}

// NewLockFreeStack is a constructor
// for the LockFreeStack type
func NewLockFreeStack() *LockFreeStack {
	newStack := new(LockFreeStack)
	newStack.Top = nil
	return newStack
}

// LockStack
type LockStack struct {
	Top        *Item
	accessLock sync.RWMutex
}

// NewLockStack is a constructor
// for the LockStack type
func NewLockStack() *LockStack {
	newStack := new(LockStack)
	newStack.Top = nil
	return newStack
}
