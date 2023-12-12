package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ilya-rusyanov/shrinklator/internal/entities"
	"github.com/ilya-rusyanov/shrinklator/internal/handlers/mocks"
	"github.com/ilya-rusyanov/shrinklator/internal/storage"
	"go.uber.org/mock/gomock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpandHandler(t *testing.T) {
	type want struct {
		code             int
		redirectLocation string
	}
	tests := []struct {
		testName string
		expect   func(*mocks.MockShrinker)
		arg      string
		want     want
	}{
		{
			testName: "nonexistent request",
			arg:      "a0",
			expect: func(m *mocks.MockShrinker) {
				m.EXPECT().
					Expand(gomock.Any(), gomock.Any()).
					Return(entities.ExpandResult{}, errors.New("err"))
			},
			want: want{
				code:             http.StatusBadRequest,
				redirectLocation: "",
			},
		},
		{
			testName: "removed URL",
			expect: func(m *mocks.MockShrinker) {
				m.EXPECT().
					Expand(gomock.Any(), gomock.Any()).
					Return(entities.ExpandResult{Removed: true}, nil)
			},
			want: want{
				code:             http.StatusGone,
				redirectLocation: "",
			},
		},
		{
			testName: "existing url",
			arg:      "664b8054bac1af66baafa7a01acd15ee",
			expect: func(m *mocks.MockShrinker) {
				m.EXPECT().
					Expand(gomock.Any(), "664b8054bac1af66baafa7a01acd15ee").
					Return(entities.ExpandResult{
						URL: "http://yandex.ru",
					},
						nil,
					)
			},
			want: want{
				code:             http.StatusTemporaryRedirect,
				redirectLocation: "http://yandex.ru",
			},
		},
	}

	var someUser *entities.UserID

	noLog := &dummyLogger{}
	storage := storage.NewInMemory(noLog)
	e := storage.Put(context.TODO(),
		"664b8054bac1af66baafa7a01acd15ee", "http://yandex.ru", someUser)
	require.NoError(t, e)

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockShrinker(ctrl)

			tc.expect(service)

			handler := NewExpand(service)

			req := httptest.NewRequest(
				http.MethodGet,
				"/"+tc.arg,
				nil)

			resp := httptest.NewRecorder()
			handler.Handler(resp, req)

			assert.Equal(t,
				tc.want.redirectLocation,
				resp.Header().Get("Location"))

			assert.Equal(t, tc.want.code, resp.Code)
		})
	}
}
