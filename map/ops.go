package lfmap

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

func (m *Map) Insert(key int64, value interface{}) error {
	newNode := NewNode(value)
	if key >= m.Size {
		return fmt.Errorf("key out of bounds")
	}
	for {
		nodeAtKey := (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&m.Items[key]))))
		oldPointer := (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&nodeAtKey))))
		newPointer := (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&newNode))))

		if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&m.Items[key])), unsafe.Pointer(oldPointer), unsafe.Pointer(newPointer)) {
			break
		}
	}
	return nil
}

// cmp compares and return which value is preferred ( Can be min, max, etc... )
func (m *Map) InsertCompare(key int64, value interface{}, cmp func(a, b interface{}) bool) (bool, error) {
	newNode := NewNode(value)
	if key >= m.Size {
		return false, fmt.Errorf("key out of bounds")
	}
	for {
		nodeAtKey := (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&m.Items[key]))))
		if nodeAtKey != nil {
			// If new exiting value has higher priotrity according to the cmp, return false
			if cmp(nodeAtKey.Value, newNode.Value) {
				return false, nil
			}
		}
		oldPointer := (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&nodeAtKey))))
		newPointer := (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&newNode))))

		if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&m.Items[key])), unsafe.Pointer(oldPointer), unsafe.Pointer(newPointer)) {
			break
		}
	}
	return true, nil
}

func (m *Map) InsertIfDoesntExist(key int64, value interface{}) (bool, error) {
	newNode := NewNode(value)
	if key >= m.Size {
		return false, fmt.Errorf("key out of bounds")
	}
	for {
		nodeAtKey := (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&m.Items[key]))))
		if nodeAtKey != nil {
			return false, nil
		}
		oldPointer := (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&nodeAtKey))))
		newPointer := (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&newNode))))

		if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&m.Items[key])), unsafe.Pointer(oldPointer), unsafe.Pointer(newPointer)) {
			break
		}
	}
	return true, nil
}

func (m *Map) Lookup(key int64) (interface{}, error) {
	if key >= m.Size {
		return nil, fmt.Errorf("key out of bounds")
	}
	for {
		nodeAtKey := (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&m.Items[key]))))

		if nodeAtKey == nil {
			return nil, fmt.Errorf("item not found")
		}

		nodePointer := (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&nodeAtKey))))
		if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&m.Items[key])), unsafe.Pointer(nodePointer), unsafe.Pointer(nodePointer)) {
			return nodePointer.Value, nil
		}
	}
}

func (m *Map) Exists(key int64) bool {
	if key >= m.Size {
		return false
	}
	for {
		nodeAtKey := (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&m.Items[key]))))
		if nodeAtKey == nil {
			return false
		}

		nodePointer := (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&nodeAtKey))))
		if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&m.Items[key])), unsafe.Pointer(nodePointer), unsafe.Pointer(nodePointer)) {
			return true
		}
	}
}

// func main() {
// 	m := NewMap(100)
// 	m.Insert(10, 2)
// 	m.Insert(25, 20)
// 	m.Insert(27, 1700)
// 	m.Insert(13, 356)
// 	m.Insert(12, 12)
// 	m.Insert(4, 23)

// 	fmt.Println(m.Lookup(25))
// 	fmt.Println(m.Lookup(13))
// 	fmt.Println(m.Lookup(2))
// 	fmt.Println(m.Lookup(105))
// }
