package services

import (
	"context"
	"fmt"

	"github.com/ilya-rusyanov/shrinklator/internal/entities"
)

type BatchStorage interface {
	PutBatch(context.Context, []entities.ShortLongPair) error
}

type Batch struct {
	storage BatchStorage
	algo    Algo
}

func NewBatch(storage BatchStorage, algorithm Algo) *Batch {
	return &Batch{
		storage: storage,
		algo:    algorithm,
	}
}

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
