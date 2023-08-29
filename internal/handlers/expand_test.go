package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ilya-rusyanov/shrinklator/internal/services"
	"github.com/ilya-rusyanov/shrinklator/internal/storage"
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
		arg      string
		want     want
	}{
		{
			testName: "empty request",
			arg:      "",
			want: want{
				code:             http.StatusBadRequest,
				redirectLocation: "",
			},
		},
		{
			testName: "nonexistent request",
			arg:      "a0",
			want: want{
				code:             http.StatusBadRequest,
				redirectLocation: "",
			},
		},
		{
			testName: "existing url",
			arg:      "664b8054bac1af66baafa7a01acd15ee",
			want: want{
				code:             http.StatusTemporaryRedirect,
				redirectLocation: "http://yandex.ru",
			},
		},
	}

	values := make(map[string]string)
	storage := storage.NewInMemory(values)
	storage.Put("664b8054bac1af66baafa7a01acd15ee", "http://yandex.ru")

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			model := services.NewShortener(storage)

			handler := NewExpand(model)

			req, err := http.NewRequest(
				http.MethodGet,
				"/"+test.arg,
				nil)

			require.NoError(t, err)

			resp := httptest.NewRecorder()
			handler.Handler().ServeHTTP(resp, req)

			assert.Equal(t,
				test.want.redirectLocation,
				resp.Header().Get("Location"))

			assert.Equal(t, test.want.code, resp.Code)
		})
	}
}
