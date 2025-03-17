package safe_cache

import (
	"fmt"
	"sync"
)

const (
	// Error
	ErrorSafeCacheIsNullPoint = "SafeCache is null point"
	ErrorSafeCacheIsEmpty     = "SafeCache length is 0"
)

type node struct {
	data   interface{}
	pNext  *node
	nodeId int
}

// SafeCache 是一个线程安全的无头单向链表，每次Set后Get能获取到最新的值。并且保证之前取出的元素还能继续使用。
// 得益于GC，SafeCache去回收那些不再被引用的对象
type SafeCache struct {
	pEnd        *node
	nodeIdCount int
	rwLock      sync.RWMutex
}

func NewSafeCache() *SafeCache {
	return &SafeCache{
		pEnd:        nil,
		nodeIdCount: 0,
	}
}

func (s *SafeCache) Set(data interface{}) error {
	if s == nil {
		return fmt.Errorf("%s", ErrorSafeCacheIsNullPoint)
	}

	s.rwLock.Lock()
	if s.pEnd == nil {
		s.pEnd = &node{
			data:   data,
			pNext:  nil,
			nodeId: s.nodeIdCount,
		}
	} else {
		s.pEnd.pNext = &node{
			data:   data,
			pNext:  nil,
			nodeId: s.nodeIdCount,
		}
		s.pEnd = s.pEnd.pNext
	}
	s.nodeIdCount++
	s.rwLock.Unlock()
	return nil
}

func (s *SafeCache) Get() (interface{}, error) {
	if s == nil {
		var r interface{}
		return r, fmt.Errorf("%s", ErrorSafeCacheIsNullPoint)
	}

	s.rwLock.RLock()
	if s.pEnd == nil {
		var r interface{}
		s.rwLock.RUnlock()
		return r, fmt.Errorf("%s", ErrorSafeCacheIsNullPoint)
	}
	res := s.pEnd.data
	s.rwLock.RUnlock()
	return res, nil
}
