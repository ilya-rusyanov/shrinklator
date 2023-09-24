package services

import (
	"context"
	"fmt"

	"github.com/ilya-rusyanov/shrinklator/internal/entities"
	"github.com/ilya-rusyanov/shrinklator/internal/logger"
	"go.uber.org/zap"
)

type ShortStorage interface {
	Put(ctx context.Context, id string, value string, uid *entities.UserID) error
	ByID(ctx context.Context, id string) (entities.ExpandResult, error)
}

type Shortener struct {
	storage ShortStorage
	algo    Algo
	log     *logger.Log
}

func NewShortener(log *logger.Log, storage ShortStorage, algorithm Algo) *Shortener {
	res := &Shortener{
		storage: storage,
		algo:    algorithm,
		log:     log,
	}
	return res
}

func (s *Shortener) Shrink(ctx context.Context, input string, uid *entities.UserID) (string, error) {
	short := s.algo(input)
	err := s.storage.Put(ctx, short, input, uid)
	s.log.Info("shortened", zap.String("short", short), zap.String("long", input))
	if err != nil {
		return "", fmt.Errorf("error storing: %w", err)
	}
	return short, nil
}

func (s *Shortener) Expand(ctx context.Context, input string) (entities.ExpandResult, error) {
	res, err := s.storage.ByID(ctx, input)

	if err != nil {
		return entities.ExpandResult{},
			fmt.Errorf("error searching: %w", err)
	}

	return res, nil
}
