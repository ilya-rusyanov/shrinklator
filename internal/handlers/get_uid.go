package handlers

import (
	"context"

	"github.com/ilya-rusyanov/shrinklator/internal/entities"
)

func getUID(ctx context.Context) *entities.UserID {
	if id := ctx.Value(UID); id != nil {
		val := id.(entities.UserID)
		return &val
	}

	return nil
}
