package lfmap

import "sync"

type Node struct {
	Value interface{}
	_pad  [6]int
}

type LockFreeMap struct {
	Size  int64
	Items []*Node
	_pad  [4]int
}

type LockMap struct {
	Size       int64
	Items      []*Node
	AccessLock sync.RWMutex
}

func NewNode(value interface{}) *Node {
	newNode := new(Node)
	newNode.Value = value
	return newNode
}

func NewLockFreeMap(size int64) *LockFreeMap {
	newMap := new(LockFreeMap)
	newMap.Size = size
	newMap.Items = make([]*Node, size)
	return newMap
}

func NewLockMap(size int64) *LockMap {
	newMap := new(LockMap)
	newMap.Size = size
	newMap.Items = make([]*Node, size)
	return newMap
}
