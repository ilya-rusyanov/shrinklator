package services

import (
	"context"
	"testing"

	"github.com/ilya-rusyanov/shrinklator/internal/entities"
	"github.com/ilya-rusyanov/shrinklator/internal/services/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestBatch(t *testing.T) {
	type want struct {
		result []entities.BatchResponse
		err    error
		dbArg  []entities.ShortLongPair
	}

	tests := []struct {
		name  string
		input []entities.BatchRequest
		want  want
	}{
		{
			name: "basic test",
			input: []entities.BatchRequest{
				{
					ID:      "a1b1c1",
					LongURL: "http://google.com",
				},
			},
			want: want{
				result: []entities.BatchResponse{
					{
						ID:       "a1b1c1",
						ShortURL: "c7b920f57e553df2bb68272f61570210",
					},
				},
				err: nil,
				dbArg: []entities.ShortLongPair{
					{
						Short: "c7b920f57e553df2bb68272f61570210",
						Long:  "http://google.com",
					},
				},
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			db := mocks.NewMockBatchStorage(ctrl)

			db.EXPECT().PutBatch(gomock.Any(),
				testCase.want.dbArg).Return(nil)

			batch := NewBatch(db, MD5Algo)

			result, err := batch.BatchShorten(context.TODO(), testCase.input)
			require.NoError(t, err)

			assert.Equal(t, testCase.want.result, result)
		})
	}
}
