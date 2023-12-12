package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testHandler struct {
	served bool
}

func (h *testHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.served = true
}

func TestPseudoAuth(t *testing.T) {
	tests := []struct {
		name string
		next bool
	}{
		{
			name: "builds cookie for client",
			next: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			auth := NewPseudoAuth(&dummyLogger{}, "key", "access_token")

			th := testHandler{}
			h := auth.Middleware(&th)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			h.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tc.next, th.served)

			var cookie *http.Cookie
			for _, c := range res.Cookies() {
				if c.Name == "access_token" {
					cookie = c
				}
			}
			assert.NotNil(t, cookie)
		})
	}
}
