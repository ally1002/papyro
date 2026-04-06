package keyring

import (
	"testing"

	"github.com/99designs/keyring"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeyring_Save(t *testing.T) {
	kr := &Keyring{ring: keyring.NewArrayKeyring(nil)}

	err := kr.Save("test-account", "secret123")
	require.NoError(t, err)

	item, err := kr.Get("test-account")
	require.NoError(t, err)

	assert.Equal(t, item.Key, "test-account")
	assert.Equal(t, item.Data, []byte("secret123"))
}
