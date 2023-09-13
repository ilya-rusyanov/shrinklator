package services

import (
	"context"
	"testing"

	"github.com/ilya-rusyanov/shrinklator/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestShortener(t *testing.T) {
	noLog := zap.NewNop()

	algo := MD5Algo

	t.Run("short new", func(t *testing.T) {
		s := NewShortener(storage.NewInMemory(noLog), algo)

		got, err := s.Shrink(context.TODO(), "http://yandex.ru")
		require.NoError(t, err)

		want := "664b8054bac1af66baafa7a01acd15ee"

		assert.Equal(t, want, got)
	})

	t.Run("expand unknown", func(t *testing.T) {
		s := NewShortener(storage.NewInMemory(noLog), algo)

		_, err := s.Expand(context.TODO(), "a")

		if err == nil {
			t.Fatal("must raise error")
		}
	})

	t.Run("expand known", func(t *testing.T) {
		s := NewShortener(storage.NewInMemory(noLog), algo)

		url := "http://yandex.ru"

		short, err := s.Shrink(context.TODO(), url)
		require.NoError(t, err)

		got, err := s.Expand(context.TODO(), short)
		require.NoError(t, err)

		assert.Equal(t, url, got)
	})

	t.Run("expand unknown", func(t *testing.T) {
		s := NewShortener(storage.NewInMemory(noLog), algo)

		_, err := s.Expand(context.TODO(), "http://google.com")

		require.Error(t, err)
	})
}
