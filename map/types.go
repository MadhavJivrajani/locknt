package lfmap

type Node struct {
	Value interface{}
}

type Map struct {
	Size  int64
	Items []*Node
}

func NewNode(value interface{}) *Node {
	newNode := new(Node)
	newNode.Value = value
	return newNode
}

func NewMap(size int64) *Map {
	newMap := new(Map)
	newMap.Size = size
	newMap.Items = make([]*Node, size)
	return newMap
}
