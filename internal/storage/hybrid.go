package storage

import (
	"fmt"
)

type Persistence interface {
	Append(short, long string) error
}

type Hybrid struct {
	memory      *inMemory
	persistence Persistence
}

func NewHybrid(inMemory *inMemory, persistence Persistence) *Hybrid {
	return &Hybrid{
		memory:      inMemory,
		persistence: persistence,
	}
}

func (s *Hybrid) Put(id, value string) error {
	s.memory.Put(id, value)

	err := s.persistence.Append(id, value)
	if err != nil {
		return fmt.Errorf("error writing to disk: %w", err)
	}

	return nil
}

func (s *Hybrid) ByID(id string) (string, error) {
	return s.memory.ByID(id)
}
