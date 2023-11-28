package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGzipCompression(t *testing.T) {
	successReply := `d8e8fca2dc0f896fd7cb4cb0031ba249`
	requestBody := `http://yandex.ru`

	t.Run("accepts gzip", func(t *testing.T) {
		received := ""

		webhook := http.HandlerFunc(
			func(rw http.ResponseWriter, r *http.Request) {
				t.Helper()

				buf, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				received = string(buf)

				_, err = io.WriteString(rw, successReply)
				require.NoError(t, err)
			})

		handler := NewGzip(&dummyLogger{})

		srv := httptest.NewServer(handler.Middleware(webhook))
		defer srv.Close()

		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(requestBody))
		require.NoError(t, err)
		err = zb.Close()
		require.NoError(t, err)

		r := httptest.NewRequest("POST", srv.URL, buf)
		r.RequestURI = ""
		r.Header.Set("Content-Encoding", "gzip")

		resp, err := http.DefaultClient.Do(r)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		defer func() {
			e := resp.Body.Close()
			require.NoError(t, e)
		}()

		assert.Equal(t, requestBody, received)
	})

	t.Run("replies with gzip", func(t *testing.T) {
		webhook := http.HandlerFunc(
			func(rw http.ResponseWriter, r *http.Request) {
				t.Helper()

				_, err := io.WriteString(rw, successReply)
				require.NoError(t, err)
			})

		handler := NewGzip(&dummyLogger{})

		srv := httptest.NewServer(handler.Middleware(webhook))
		defer srv.Close()

		buf := bytes.NewBufferString(requestBody)
		r := httptest.NewRequest("POST", srv.URL, buf)
		r.RequestURI = ""
		r.Header.Set("Accept-Encoding", "gzip")

		resp, err := http.DefaultClient.Do(r)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		defer func() {
			e := resp.Body.Close()
			require.NoError(t, e)
		}()

		zr, err := gzip.NewReader(resp.Body)
		require.NoError(t, err)

		b, err := io.ReadAll(zr)
		require.NoError(t, err)

		require.Equal(t, successReply, string(b))
	})
}
