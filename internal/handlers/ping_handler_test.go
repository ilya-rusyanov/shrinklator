package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/ilya-rusyanov/shrinklator/internal/handlers/mocks"
	"github.com/stretchr/testify/assert"
)

func TestPingHandler(t *testing.T) {
	tests := []struct {
		name        string
		expectation func(*mocks.MockPinger)
		wantStatus  int
	}{
		{
			name: "ping ok",
			expectation: func(m *mocks.MockPinger) {
				m.EXPECT().
					Ping(gomock.Any()).
					Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "ping fail",
			expectation: func(m *mocks.MockPinger) {
				m.EXPECT().
					Ping(gomock.Any()).
					Return(errors.New("abstract error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pinger := mocks.NewMockPinger(ctrl)

			tc.expectation(pinger)

			handler := NewPing(&dummyLogger{}, pinger)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			handler.Handler(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tc.wantStatus, res.StatusCode)
		})
	}
}
