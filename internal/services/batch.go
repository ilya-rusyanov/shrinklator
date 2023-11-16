package services

import (
	"context"
	"fmt"

	"github.com/ilya-rusyanov/shrinklator/internal/entities"
)

// BatchStorage - storage for URLs shortened in batch requests
type BatchStorage interface {
	PutBatch(context.Context, []entities.ShortLongPair) error
}

// Batch - usecase for shortening URLs in batch
type Batch struct {
	storage BatchStorage
	algo    Algo
}

// NewBatch constructs Batch objects
func NewBatch(storage BatchStorage, algorithm Algo) *Batch {
	return &Batch{
		storage: storage,
		algo:    algorithm,
	}
}

// BatchShorten bulk shorten URLs
func (b *Batch) BatchShorten(ctx context.Context,
	input []entities.BatchRequest) ([]entities.BatchResponse, error) {
	toStore := make([]entities.ShortLongPair, len(input))
	result := make([]entities.BatchResponse, len(input))

	for i, val := range input {
		toStore[i].Short = b.algo(val.LongURL)
		toStore[i].Long = val.LongURL
		result[i].ID = val.ID
		result[i].ShortURL = toStore[i].Short
	}

	err := b.storage.PutBatch(ctx, toStore)

	if err != nil {
		return nil, fmt.Errorf("failed to store batch: %w", err)
	}

	return result, nil
}
