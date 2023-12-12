package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/ilya-rusyanov/shrinklator/internal/entities"
	"github.com/ilya-rusyanov/shrinklator/internal/handlers/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestShortenHandler(t *testing.T) {
	type input struct {
		body string
		uid  *entities.UserID
	}

	type want struct {
		code int
		body string
	}

	someUID := func() *entities.UserID {
		id := entities.UserID("a")
		return &id
	}

	tests := []struct {
		testName string
		input    input
		expect   func(*mocks.MockShrinker)
		want     want
	}{
		{
			testName: "service conflict error",
			input: input{
				body: "http://yandex.ru",
			},
			expect: func(m *mocks.MockShrinker) {
				m.EXPECT().
					Shrink(gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", errors.New("generic error"))
			},
			want: want{
				code: http.StatusBadRequest,
				body: "",
			},
		},
		{
			testName: "successfull url shortening",
			input: input{
				body: "http://yandex.ru",
				uid:  someUID(),
			},
			expect: func(m *mocks.MockShrinker) {
				m.EXPECT().
					Shrink(gomock.Any(), "http://yandex.ru", someUID()).
					Return("short", nil)
			},
			want: want{
				code: http.StatusCreated,
				body: "short",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockShrinker(ctrl)

			tc.expect(service)

			noLog := dummyLogger{}

			shortenHandler := NewShorten(
				&noLog,
				service,
				"http://localhost:8080",
			)

			req := httptest.NewRequest(
				http.MethodPost,
				"/",
				strings.NewReader(tc.input.body))

			if tc.input.uid != nil {
				ctx := context.WithValue(
					context.Background(),
					UID,
					*tc.input.uid)
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()

			shortenHandler.Handler(w, req)

			res := w.Result()
			defer res.Body.Close()

			respBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, tc.want.code, res.StatusCode)

			if res.StatusCode == http.StatusCreated {
				assert.Regexp(
					t,
					regexp.MustCompile("http://.+/"+tc.want.body),
					string(respBody))
			}
		})
	}
}
