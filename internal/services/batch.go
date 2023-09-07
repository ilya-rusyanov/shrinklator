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
	return nil, fmt.Errorf("TODO")
}
