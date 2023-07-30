package shortener

import (
	"testing"
)

func TestShortener(t *testing.T) {
	t.Run("short new", func(t *testing.T) {
		s := New()

		got, err := s.Shrink("http://yandex.ru")

		if err != nil {
			t.Fatal("got error")
		}

		want := "0"

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("expand unknown", func(t *testing.T) {
		s := New()

		_, err := s.Expand("a")

		if err == nil {
			t.Fatal("must raise error")
		}
	})

	t.Run("expand known", func(t *testing.T) {
		s := New()

		url := "http://yandex.ru"

		short, err := s.Shrink(url)

		if err != nil {
			t.Fatal("must be valid operation")
		}

		got, err := s.Expand(short)

		if err != nil {
			t.Fatal("must be valid operation")
		}

		if got != url {
			t.Errorf("got %q want %q", got, url)
		}
	})

	t.Run("expand unknown", func(t *testing.T) {
		s := New()

		_, err := s.Expand("http://google.com")

		if err == nil {
			t.Errorf("want error, got nil")
		}
	})
}
