package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ilya-rusyanov/shrinklator/internal/services"
	"github.com/ilya-rusyanov/shrinklator/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShortenRESThandler(t *testing.T) {
	type want struct {
		code         int
		responseBody string
	}

	tests := []struct {
		name  string
		input string
		want  want
	}{
		{
			name:  "empty invalid",
			input: "",
			want: want{
				code:         http.StatusBadRequest,
				responseBody: "unexpected end of JSON input\n",
			},
		},
		{
			name:  "valid",
			input: `{ "url": "https://practicum.yandex.ru" }`,
			want: want{
				code:         http.StatusCreated,
				responseBody: `{"result":"http://localhost/6bdb5b0e26a76e4dab7cd1a272caebc0"}`,
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			values := make(map[string]string)
			storage := storage.NewInMemory(values)
			service := services.NewShortener(storage)
			handler := NewShortenREST(service, "http://localhost")

			req, err := http.NewRequest(
				http.MethodPost,
				"/",
				strings.NewReader(testCase.input))
			require.NoError(t, err)

			resp := httptest.NewRecorder()
			handler.Handler().ServeHTTP(resp, req)

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, testCase.want.code, resp.Code)
			assert.Equal(t, testCase.want.responseBody, string(respBody))
		})
	}
}
