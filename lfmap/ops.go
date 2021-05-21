package lfmap

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

func (m *LockFreeMap) Insert(key int64, value interface{}) error {
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
func (m *LockFreeMap) InsertCompare(key int64, value interface{}, cmp func(a, b interface{}) bool) (bool, error) {
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

func (m *LockFreeMap) InsertIfDoesntExist(key int64, value interface{}) (bool, error) {
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

func (m *LockFreeMap) Lookup(key int64) (interface{}, error) {
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

func (m *LockFreeMap) Exists(key int64) bool {
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

func (m *LockMap) Insert(key int64, value interface{}) error {
	m.AccessLock.Lock()
	defer m.AccessLock.Unlock()
	newNode := NewNode(value)
	if key >= m.Size {
		return fmt.Errorf("key out of bounds")
	}
	m.Items[key] = newNode
	return nil
}

func (m *LockMap) Lookup(key int64) (interface{}, error) {
	m.AccessLock.RLock()
	defer m.AccessLock.RUnlock()
	if key >= m.Size {
		return nil, fmt.Errorf("key out of bounds")
	}
	nodeAtKey := m.Items[key]
	if nodeAtKey == nil {
		return nil, fmt.Errorf("key not found")
	}

	return nodeAtKey.Value, nil
}
