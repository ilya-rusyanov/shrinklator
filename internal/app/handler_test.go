package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ilya-rusyanov/shrinklator/internal/models"
	"github.com/stretchr/testify/assert"
)

type fakeShrinker struct {
}

func (s *fakeShrinker) Shrink(string) (string, error) {
	return "", nil
}

func (s *fakeShrinker) Expand(string) (string, error) {
	return "http://yandex.ru", nil
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
			h := shortenerHandler{models.New()}

			request := httptest.NewRequest(test.method, test.path, nil)

			w := httptest.NewRecorder()
			h.ServeHTTP(w, request)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want.code, res.StatusCode)
		})
	}

	t.Run("post", func(t *testing.T) {
		h := shortenerHandler{models.New()}

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
