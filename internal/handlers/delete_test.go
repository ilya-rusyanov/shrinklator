package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	"github.com/ilya-rusyanov/shrinklator/internal/entities"
	"github.com/ilya-rusyanov/shrinklator/internal/handlers/mocks"
)

func TestDelete(t *testing.T) {
	type input struct {
		body string
		user entities.UserID
	}

	tests := []struct {
		name  string
		input input
		want  int
	}{
		{
			name: "broken json",
			input: input{
				body: "{",
			},
			want: http.StatusBadRequest,
		},
		{
			name: "unauthorized user",
			input: input{
				body: "null",
			},
			want: http.StatusUnauthorized,
		},
		{
			name: "valid user",
			input: input{
				body: "null",
				user: entities.UserID("2"),
			},
			want: http.StatusAccepted,
		},
		{
			name: "successfull delete",
			input: input{
				body: `["a", "b", "c"]`,
				user: entities.UserID("2"),
			},
			want: http.StatusAccepted,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockDeleteService(ctrl)
			service.EXPECT().
				Delete(gomock.Any(), gomock.Any()).
				AnyTimes()

			handler := NewDeleteHandler(&dummyLogger{}, service)

			req, err := http.NewRequest(
				http.MethodDelete,
				"/api/user/urls",
				strings.NewReader(tc.input.body),
			)
			require.NoError(t, err)

			if tc.input.user != "" {
				ctx := context.WithValue(
					context.Background(),
					UID,
					tc.input.user,
				)
				req = req.WithContext(ctx)
			}

			resp := httptest.NewRecorder()
			handler.Handler(resp, req)

			assert.Equal(t, tc.want, resp.Code)
		})
	}
}
