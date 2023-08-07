package storage

import (
	"errors"
	"sync"
)

var errNotFound = errors.New("not found")

type inMemory struct {
	data  map[string]string
	mutex sync.Mutex
}

func NewInMemory() *inMemory {
	res := inMemory{}
	res.data = map[string]string{}
	return &res
}

func (s *inMemory) Put(id, value string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.data[id] = value
}

func (s *inMemory) ByID(id string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	value, ok := s.data[id]

	if !ok {
		return "", errNotFound
	}

	return value, nil
}
