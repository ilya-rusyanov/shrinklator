package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/ilya-rusyanov/shrinklator/internal/services"
	"github.com/ilya-rusyanov/shrinklator/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShortenHandler(t *testing.T) {
	type want struct {
		code     int
		response string
	}

	tests := []struct {
		testName string
		body     string
		want     want
	}{
		{
			testName: "url shortening",
			body:     "http://yandex.ru",
			want: want{
				code:     http.StatusCreated,
				response: "664b8054bac1af66baafa7a01acd15ee",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			noLog := dummyLogger{}
			storage := storage.NewInMemory(&noLog)
			model := services.NewShortener(&noLog, storage, services.MD5Algo)

			shortenHandler := NewShorten(&noLog, model,
				"http://localhost:8080")

			server := httptest.NewServer(http.HandlerFunc(shortenHandler.Handler))
			defer server.Close()

			req, err := http.NewRequest(
				http.MethodPost,
				server.URL+"/",
				strings.NewReader(test.body))

			require.NoError(t, err)

			resp, err := server.Client().Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Regexp(t,
				regexp.MustCompile("http://.+/"+test.want.response),
				string(respBody))

			assert.Equal(t, test.want.code, resp.StatusCode)
		})
	}
}
