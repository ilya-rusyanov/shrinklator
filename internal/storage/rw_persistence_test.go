package storage

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRWPersistence(t *testing.T) {
	t.Run("appending", func(t *testing.T) {
		buf := bytes.Buffer{}
		p := NewRWPersistence(&buf, &buf)
		p.Append("4rSPg8ap", "http://yandex.ru")
		p.Append("edVPg3ks", "http://ya.ru")
		assert.Equal(t,
			`{"uuid":"1","short_url":"4rSPg8ap","original_url":"http://yandex.ru"}
{"uuid":"2","short_url":"edVPg3ks","original_url":"http://ya.ru"}
`,
			buf.String())
	})

	t.Run("reading", func(t *testing.T) {
		buf := bytes.Buffer{}
		_, err := buf.Write([]byte(
			`{"uuid":"1","short_url":"4rSPg8ap","original_url":"http://yandex.ru"}
{"uuid":"2","short_url":"edVPg3ks","original_url":"http://ya.ru"}
{"uuid":"3","short_url":"dG56Hqxm","original_url":"http://practicum.yandex.ru"}
`))
		require.NoError(t, err)

		p := NewRWPersistence(&buf, &buf)
		values, err := p.ReadAll()
		require.NoError(t, err)

		require.Equal(t, 3, len(values))
		assert.Equal(t, "http://yandex.ru", values["4rSPg8ap"])
	})
}
