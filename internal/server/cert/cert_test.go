package cert

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	s := http.Server{}

	err := New()(&s)
	require.NoError(t, err)

	assert.NotNil(t, s.TLSConfig)
}
