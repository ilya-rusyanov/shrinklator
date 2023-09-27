package storage

import (
	"context"

	"github.com/ilya-rusyanov/shrinklator/internal/entities"
)

type Interface interface {
	Put(ctx context.Context, id, value string, uid *entities.UserID) error
	PutBatch(context.Context, []entities.ShortLongPair) error
	ByID(ctx context.Context, id string) (entities.ExpandResult, error)
	ByUID(context.Context, entities.UserID) (entities.PairArray, error)
	Delete(context.Context, entities.DeleteRequest) error
	MustClose()
	Ping(context.Context) error
}
