package storage

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/ilya-rusyanov/shrinklator/internal/entities"
	"go.uber.org/zap"
)

// InMemory - storage for objects in RAM
type InMemory struct {
	data  map[string]string
	mutex sync.RWMutex
	log   Logger
}

// NewInMemory constructs InMemory objects
func NewInMemory(log Logger) *InMemory {
	return &InMemory{
		data: make(map[string]string),
		log:  log,
	}
}

// MustClose finilazes object or panics
func (s *InMemory) MustClose() {
}

// Put adds entry
func (s *InMemory) Put(ctx context.Context, id, value string, uid *entities.UserID) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.log.Info("store", zap.String("id", id),
		zap.String("value", value))

	if val, ok := s.data[id]; ok {
		return ErrAlreadyExists{
			StoredValue: val,
		}
	}

	s.data[id] = value

	return nil
}

// PutBatch adds multiple entries in bulk
func (s *InMemory) PutBatch(ctx context.Context, data []entities.ShortLongPair) error {
	return fmt.Errorf("TODO")
}

// ByID searches for entry by identifier
func (s *InMemory) ByID(ctx context.Context, id string) (entities.ExpandResult, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	value, ok := s.data[id]

	if !ok {
		s.log.Info("cannot find record", zap.String("id", id))
		return entities.ExpandResult{}, ErrNotFound
	}

	s.log.Info("successuflly found record", zap.String("id", id),
		zap.String("value", value))

	return entities.ExpandResult{URL: value}, nil
}

// ByUID searches entries by user identifier
func (s *InMemory) ByUID(context.Context, entities.UserID) (entities.PairArray, error) {
	return nil, errors.New("TODO")
}

// Delete deletes entry
func (s *InMemory) Delete(context.Context, entities.DeleteRequest) error {
	return errors.New("TODO")
}

// Ping checks storage accessibility
func (s *InMemory) Ping(context.Context) error {
	return nil
}
