package storage

import (
	"sync"
)

// Storage is default template storage
type Storage[T any] struct {
	data map[string]T
	mu   *sync.Mutex
}

func NewStorage[T any]() *Storage[T] {
	return &Storage[T]{
		data: make(map[string]T),
		mu:   &sync.Mutex{},
	}
}

func (s *Storage[T]) Set(key string, val T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = val
}

func (s *Storage[T]) Get(key string) (val T, exists bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	val, exists = s.data[key]
	return
}
