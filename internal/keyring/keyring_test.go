package keyring

import (
	"testing"

	"github.com/99designs/keyring"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveAndGet(t *testing.T) {
	kr := &Keyring{ring: keyring.NewArrayKeyring(nil)}

	err := kr.Save("test-account", "secret123")
	require.NoError(t, err)

	item, err := kr.Get("test-account")
	require.NoError(t, err)

	assert.Equal(t, item.Key, "test-account")
	assert.Equal(t, item.Data, []byte("secret123"))
}

func TestSaveAlreadyExistent(t *testing.T) {
	kr := &Keyring{ring: keyring.NewArrayKeyring(nil)}

	err := kr.Save("test-account", "secret123")
	require.NoError(t, err)

	err = kr.Save("test-account", "secret321")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "password for account \"test-account\" already exists")
}

func TestGetNotFound(t *testing.T) {
	kr := &Keyring{ring: keyring.NewArrayKeyring(nil)}

	_, err := kr.Get("non-existent")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "The specified item could not be found in the keyring")
}

func TestDelete(t *testing.T) {
	kr := &Keyring{ring: keyring.NewArrayKeyring(nil)}

	err := kr.Save("test-account", "secret123")
	require.NoError(t, err)

	err = kr.Delete("test-account")
	require.NoError(t, err)

	_, err = kr.Get("test-account")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "The specified item could not be found in the keyring")
}
