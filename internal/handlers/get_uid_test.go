package handlers

import (
	"context"
	"testing"

	"github.com/ilya-rusyanov/shrinklator/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestGetUID(t *testing.T) {
	ptrID := func(id string) *entities.UserID {
		res := entities.UserID(id)
		return &res
	}

	tests := []struct {
		name string
		val  *entities.UserID
		want *entities.UserID
	}{
		{
			name: "has value",
			val:  ptrID("one"),
			want: ptrID("one"),
		},
		{
			name: "no value",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			if tc.val != nil {
				ctx = context.WithValue(ctx, UID, *tc.val)
			}

			got := getUID(ctx)

			assert.Equal(t, tc.want, got)
		})
	}
}
