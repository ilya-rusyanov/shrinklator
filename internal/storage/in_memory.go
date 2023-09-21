package storage

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/ilya-rusyanov/shrinklator/internal/entities"
	"github.com/ilya-rusyanov/shrinklator/internal/logger"
	"go.uber.org/zap"
)

type InMemory struct {
	data  map[string]string
	mutex sync.RWMutex
	log   *logger.Log
}

func NewInMemory(log *logger.Log) *InMemory {
	return &InMemory{
		data: make(map[string]string),
		log:  log,
	}
}

func (s *InMemory) MustClose() {
}

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

func (s *InMemory) PutBatch(ctx context.Context, data []entities.ShortLongPair) error {
	return fmt.Errorf("TODO")
}

func (s *InMemory) ByID(ctx context.Context, id string) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	value, ok := s.data[id]

	if !ok {
		s.log.Info("cannot find record", zap.String("id", id))
		return "", ErrNotFound
	}

	s.log.Info("successuflly found record", zap.String("id", id),
		zap.String("value", value))

	return value, nil
}

func (s *InMemory) ByUID(context.Context, entities.UserID) (entities.PairArray, error) {
	return nil, errors.New("TODO")
}

func (s *InMemory) Ping(context.Context) error {
	return nil
}
