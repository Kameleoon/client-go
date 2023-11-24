package storage

import "sync"

type DataCollectionStorage[V any] interface {
	Enumerate(f func(V) bool)
	Last() V
	Len() int
}

func enumerateSlice[V any](s []V, f func(V) bool) bool {
	for _, v := range s {
		if !f(v) {
			return false
		}
	}
	return true
}

type DataCollectionStorageImpl[V any] struct {
	mx   *sync.RWMutex
	data []V
}

func NewDataCollectionStorageImpl[V any](mx *sync.RWMutex, data []V) *DataCollectionStorageImpl[V] {
	return &DataCollectionStorageImpl[V]{
		mx:   mx,
		data: data,
	}
}

func (s *DataCollectionStorageImpl[V]) Enumerate(f func(V) bool) {
	if s.data != nil {
		s.mx.RLock()
		defer s.mx.RUnlock()
		enumerateSlice[V](s.data, f)
	}
}

func (s *DataCollectionStorageImpl[V]) Last() V {
	if s.data != nil {
		s.mx.RLock()
		defer s.mx.RUnlock()
		i := len(s.data) - 1
		if i >= 0 {
			return s.data[i]
		}
	}
	var defaultV V
	return defaultV
}

func (s *DataCollectionStorageImpl[V]) Len() int {
	if s.data != nil {
		s.mx.RLock()
		defer s.mx.RUnlock()
		return len(s.data)
	}
	return 0
}
