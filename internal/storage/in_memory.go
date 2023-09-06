package storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/ilya-rusyanov/shrinklator/internal/entities"
	"github.com/ilya-rusyanov/shrinklator/internal/logger"
	"go.uber.org/zap"
)

type inMemory struct {
	data  map[string]string
	mutex sync.RWMutex
	log   *logger.Log
}

func NewInMemory(log *logger.Log) *inMemory {
	return &inMemory{
		data: make(map[string]string),
		log:  log,
	}
}

func (s *inMemory) Put(ctx context.Context, id, value string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.log.Info("store", zap.String("id", id),
		zap.String("value", value))

	s.data[id] = value

	return nil
}

func (s *inMemory) PutBatch(ctx context.Context, data []entities.ShortLongPair) error {
	return fmt.Errorf("TODO")
}

func (s *inMemory) ByID(ctx context.Context, id string) (string, error) {
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
