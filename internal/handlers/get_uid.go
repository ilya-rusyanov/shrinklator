package handlers

import (
	"context"

	"github.com/ilya-rusyanov/shrinklator/internal/entities"
)

func getUID(ctx context.Context) (id entities.UserID, err error) {
	if id := ctx.Value(UID); id != nil {
		return id.(entities.UserID), nil
	}

	return entities.UserID(0), errNoUserID
}
