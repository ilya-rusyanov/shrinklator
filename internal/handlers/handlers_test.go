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
			storage := storage.NewInMemory()
			model := services.NewShortener(storage)

			server := httptest.NewServer(
				Shorten(model, "http://localhost:8080"))
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

	storage := storage.NewInMemory()
	storage.Put("664b8054bac1af66baafa7a01acd15ee", "http://yandex.ru")

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			model := services.NewShortener(storage)
			handler := Expand(model)

			req, err := http.NewRequest(
				http.MethodGet,
				"/"+test.arg,
				nil)

			require.NoError(t, err)

			resp := httptest.NewRecorder()
			handler.ServeHTTP(resp, req)

			assert.Equal(t,
				test.want.redirectLocation,
				resp.Header().Get("Location"))

			assert.Equal(t, test.want.code, resp.Code)
		})
	}
}

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
			storage := storage.NewInMemory()
			service := services.NewShortener(storage)
			handler := ShortenREST(service, "http://localhost")

			req, err := http.NewRequest(
				http.MethodPost,
				"/",
				strings.NewReader(testCase.input))
			require.NoError(t, err)

			resp := httptest.NewRecorder()
			handler.ServeHTTP(resp, req)

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, testCase.want.code, resp.Code)
			assert.Equal(t, testCase.want.responseBody, string(respBody))
		})
	}
}
