package services

import (
	"context"
	"fmt"

	"github.com/ilya-rusyanov/shrinklator/internal/entities"
)

// ShortStorage - represents storage for short URLs
type ShortStorage interface {
	Put(ctx context.Context, id string, value string, uid *entities.UserID) error
	ByID(ctx context.Context, id string) (entities.ExpandResult, error)
}

// Shortener - usecase for shortening URLs
type Shortener struct {
	storage ShortStorage
	algo    Algo
	log     Logger
}

// NewShortener constructs Shortener objects
func NewShortener(log Logger, storage ShortStorage, algorithm Algo) *Shortener {
	res := &Shortener{
		storage: storage,
		algo:    algorithm,
		log:     log,
	}
	return res
}

// Shrink shortens URL
func (s *Shortener) Shrink(ctx context.Context, input string, uid *entities.UserID) (string, error) {
	short := s.algo(input)
	err := s.storage.Put(ctx, short, input, uid)
	s.log.Infof("shortened %q, long %q", short, input)
	if err != nil {
		return "", fmt.Errorf("error storing: %w", err)
	}
	return short, nil
}

// Expand expands URL
func (s *Shortener) Expand(ctx context.Context, input string) (entities.ExpandResult, error) {
	res, err := s.storage.ByID(ctx, input)

	if err != nil {
		return entities.ExpandResult{},
			fmt.Errorf("error searching: %w", err)
	}

	return res, nil
}
