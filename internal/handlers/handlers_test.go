package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/ilya-rusyanov/shrinklator/internal/models"
	"github.com/ilya-rusyanov/shrinklator/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeShrinker struct {
}

func (s *fakeShrinker) Shrink(string) string {
	return ""
}

func (s *fakeShrinker) Expand(string) (string, error) {
	return "http://yandex.ru", nil
}

func TestPostHandler(t *testing.T) {
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
			storage := storage.New()
			model := models.New(storage)

			server := httptest.NewServer(postHandler(model))
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

func TestHandler(t *testing.T) {
	type want struct {
		code int
	}
	errTests := []struct {
		name   string
		method string
		path   string
		want   want
	}{
		{
			name:   "unsupported method",
			method: http.MethodHead,
			path:   "/",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:   "unsupported path",
			method: http.MethodPost,
			path:   "/nowhere",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:   "unknown url",
			method: http.MethodGet,
			path:   "/a0",
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}

	for _, test := range errTests {
		t.Run(test.name, func(t *testing.T) {
			h := shortenerHandler{models.New(storage.New())}

			request := httptest.NewRequest(test.method, test.path, nil)

			w := httptest.NewRecorder()
			h.ServeHTTP(w, request)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want.code, res.StatusCode)
		})
	}

	t.Run("post", func(t *testing.T) {
		h := shortenerHandler{models.New(storage.New())}

		bodyReader := strings.NewReader("http://yandex.ru")

		request := httptest.NewRequest(http.MethodPost, "/", bodyReader)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, request)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusCreated, res.StatusCode)
	})

	t.Run("get", func(t *testing.T) {
		h := shortenerHandler{&fakeShrinker{}}

		request := httptest.NewRequest(http.MethodGet, "/0", nil)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, request)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusTemporaryRedirect, res.StatusCode)
	})
}
