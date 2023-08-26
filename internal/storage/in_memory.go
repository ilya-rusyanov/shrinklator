package storage

import (
	"errors"
	"sync"

	"github.com/ilya-rusyanov/shrinklator/internal/logger"
	"go.uber.org/zap"
)

var errNotFound = errors.New("not found")

type inMemory struct {
	data  map[string]string
	mutex sync.RWMutex
}

func NewInMemory(initialValues map[string]string) *inMemory {
	return &inMemory{
		data: initialValues,
	}
}

func (s *inMemory) Put(id, value string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	logger.Log.Info("store", zap.String("id", id),
		zap.String("value", value))

	s.data[id] = value

	return nil
}

func (s *inMemory) ByID(id string) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	value, ok := s.data[id]

	if !ok {
		logger.Log.Info("cannot find record", zap.String("id", id))
		return "", errNotFound
	}

	logger.Log.Info("successuflly found record", zap.String("id", id),
		zap.String("value", value))

	return value, nil
}
