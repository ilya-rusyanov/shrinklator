package storage

import (
	"context"

	"github.com/ilya-rusyanov/shrinklator/internal/entities"
)

type Interface interface {
	Put(ctx context.Context, id, value string) error
	PutBatch(context.Context, []entities.ShortLongPair) error
	ByID(ctx context.Context, id string) (string, error)
	ByUID(context.Context, entities.UserID) (entities.PairArray, error)
	MustClose()
	Ping(context.Context) error
}
