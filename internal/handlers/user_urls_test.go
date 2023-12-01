package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/ilya-rusyanov/shrinklator/internal/entities"
	"github.com/ilya-rusyanov/shrinklator/internal/handlers/mocks"
)

func TestUserURLs(t *testing.T) {
	type input struct {
		user *entities.UserID
		body string
	}

	type want struct {
		status int
		body   string
	}

	someUser := func() *entities.UserID {
		res := entities.UserID("u")
		return &res
	}

	tests := []struct {
		name        string
		input       input
		expectation func(*mocks.MockURLsService)
		want        want
	}{
		{
			name: "no user",
			want: want{
				status: http.StatusUnauthorized,
			},
		},
		{
			name: "sevice error",
			input: input{
				user: someUser(),
			},
			expectation: func(m *mocks.MockURLsService) {
				m.EXPECT().
					URLsForUser(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("arbitrary error"))
			},
			want: want{
				status: http.StatusInternalServerError,
			},
		},
		{
			name: "no urls",
			input: input{
				user: someUser(),
			},
			expectation: func(m *mocks.MockURLsService) {
				m.EXPECT().
					URLsForUser(gomock.Any(), gomock.Any())
			},
			want: want{
				status: http.StatusNoContent,
			},
		},
		{
			name: "all ok",
			input: input{
				user: someUser(),
			},
			expectation: func(m *mocks.MockURLsService) {
				m.EXPECT().
					URLsForUser(gomock.Any(), gomock.Any()).
					Return(entities.PairArray{
						entities.ShortLongPair{
							Short: "a",
							Long:  "ab",
						},
					}, nil)
			},
			want: want{
				status: http.StatusOK,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockURLsService(ctrl)

			if tc.expectation != nil {
				tc.expectation(service)
			}

			handler := NewUserURLs(&dummyLogger{}, service, "")

			req := httptest.NewRequest(http.MethodGet, "/", nil)

			ctx := context.Background()

			if tc.input.user != nil {
				ctx = context.WithValue(ctx, UID, *tc.input.user)
			}

			req = req.WithContext(ctx)

			w := httptest.NewRecorder()

			handler.Handler(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tc.want.status, res.StatusCode)
		})
	}
}
