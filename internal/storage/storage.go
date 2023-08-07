package storage

import (
	"errors"
	"sync"
)

var errNotFound = errors.New("not found")

type Storage struct {
	data  map[string]string
	mutex sync.Mutex
}

func New() *Storage {
	res := Storage{}
	res.data = map[string]string{}
	return &res
}

func (s *Storage) Put(id, value string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.data[id] = value
}

func (s *Storage) ByID(id string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	value, ok := s.data[id]

	if !ok {
		return "", errNotFound
	}

	return value, nil
}
