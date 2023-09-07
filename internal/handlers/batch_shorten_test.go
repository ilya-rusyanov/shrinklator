package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ilya-rusyanov/shrinklator/internal/entities"
	"github.com/ilya-rusyanov/shrinklator/internal/handlers/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestBatchShorten(t *testing.T) {
	type want struct {
		serviceArgument []entities.BatchRequest
		code            int
		responseBody    string
	}

	tests := []struct {
		name          string
		input         string
		serviceReturn []entities.BatchResponse
		serviceError  error
		want          want
	}{
		{
			name: "single entry",
			input: `[
    {
        "correlation_id": "a1b2c3",
	"original_url": "http://google.com"
    }
]`,
			serviceReturn: []entities.BatchResponse{
				entities.BatchResponse{
					ID:       "a1b2c3",
					ShortURL: "c7b920f57e553df2bb68272f61570210",
				},
			},
			want: want{
				serviceArgument: []entities.BatchRequest{
					entities.BatchRequest{
						ID:      "a1b2c3",
						LongURL: "http://google.com",
					},
				},
				code: http.StatusCreated,
				responseBody: `[{"correlation_id":"a1b2c3","short_url":"http://localhost/c7b920f57e553df2bb68272f61570210"}]
`,
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockBatchServicer(ctrl)
			service.EXPECT().
				BatchShorten(gomock.Any(), testCase.want.serviceArgument).
				Return(testCase.serviceReturn, testCase.serviceError)

			noLog := zap.NewNop()
			handler := NewBatchShorten(noLog, service, "http://localhost")

			req, err := http.NewRequest(
				http.MethodPost,
				"/",
				strings.NewReader(testCase.input))
			require.NoError(t, err)

			resp := httptest.NewRecorder()
			handler.Handler(resp, req)

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, testCase.want.code, resp.Code)
			assert.Equal(t, testCase.want.responseBody, string(respBody))
		})
	}
}
