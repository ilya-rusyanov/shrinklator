package models

import (
	"testing"

	"github.com/ilya-rusyanov/shrinklator/internal/storage"
)

func TestShortener(t *testing.T) {
	t.Run("short new", func(t *testing.T) {
		s := NewShortenerService(storage.NewInMemory())

		got := s.Shrink("http://yandex.ru")

		want := "664b8054bac1af66baafa7a01acd15ee"

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("expand unknown", func(t *testing.T) {
		s := NewShortenerService(storage.NewInMemory())

		_, err := s.Expand("a")

		if err == nil {
			t.Fatal("must raise error")
		}
	})

	t.Run("expand known", func(t *testing.T) {
		s := NewShortenerService(storage.NewInMemory())

		url := "http://yandex.ru"

		short := s.Shrink(url)

		got, err := s.Expand(short)

		if err != nil {
			t.Fatal("must be valid operation")
		}

		if got != url {
			t.Errorf("got %q want %q", got, url)
		}
	})

	t.Run("expand unknown", func(t *testing.T) {
		s := NewShortenerService(storage.NewInMemory())

		_, err := s.Expand("http://google.com")

		if err == nil {
			t.Errorf("want error, got nil")
		}
	})
}
