package services

import (
	"testing"

	"github.com/ilya-rusyanov/shrinklator/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestShortener(t *testing.T) {
	noLog := zap.NewNop()

	t.Run("short new", func(t *testing.T) {
		s := NewShortener(storage.NewInMemory(noLog, makeValues(t)))

		got, err := s.Shrink("http://yandex.ru")
		require.NoError(t, err)

		want := "664b8054bac1af66baafa7a01acd15ee"

		assert.Equal(t, want, got)
	})

	t.Run("expand unknown", func(t *testing.T) {
		s := NewShortener(storage.NewInMemory(noLog, makeValues(t)))

		_, err := s.Expand("a")

		if err == nil {
			t.Fatal("must raise error")
		}
	})

	t.Run("expand known", func(t *testing.T) {
		s := NewShortener(storage.NewInMemory(noLog, makeValues(t)))

		url := "http://yandex.ru"

		short, err := s.Shrink(url)
		require.NoError(t, err)

		got, err := s.Expand(short)
		require.NoError(t, err)

		assert.Equal(t, url, got)
	})

	t.Run("expand unknown", func(t *testing.T) {
		s := NewShortener(storage.NewInMemory(noLog, makeValues(t)))

		_, err := s.Expand("http://google.com")

		require.Error(t, err)
	})
}

func makeValues(tb testing.TB) map[string]string {
	return make(map[string]string)
}
