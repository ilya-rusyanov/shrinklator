package storage

import (
	"errors"
	"sync"
)

var errNotFound = errors.New("not found")

type inMemory struct {
	data  map[string]string
	mutex sync.RWMutex
}

func NewInMemory() *inMemory {
	res := inMemory{}
	res.data = make(map[string]string)
	return &res
}

func (s *inMemory) Put(id, value string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.data[id] = value
}

func (s *inMemory) ByID(id string) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	value, ok := s.data[id]

	if !ok {
		return "", errNotFound
	}

	return value, nil
}
