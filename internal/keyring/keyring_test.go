package keyring

import (
	"testing"

	"github.com/99designs/keyring"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSave(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		wantErr bool
	}{
		{"saves password successfully", "test-account-1", false},
		{"fails when password exists", "test-account-2", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kr := &Keyring{ring: keyring.NewArrayKeyring(nil)}

			err := kr.Save(tt.key, "secret123")
			require.NoError(t, err)

			if tt.wantErr {
				err = kr.Save(tt.key, "secret321")
				require.Error(t, err)
				assert.Contains(t, err.Error(), "already exists")
			}
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		wantErr  bool
		wantData string
	}{
		{"returns stored password", "test-account", false, "secret123"},
		{"fails when key is missing", "non-existent", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kr := &Keyring{ring: keyring.NewArrayKeyring(nil)}
			_ = kr.Save("test-account", "secret123")

			item, err := kr.Get(tt.key)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "could not be found")
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.key, item.Key)
				assert.Equal(t, []byte(tt.wantData), item.Data)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		wantErr bool
	}{
		{"deletes existing key", "test-account", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kr := &Keyring{ring: keyring.NewArrayKeyring(nil)}
			_ = kr.Save("test-account", "secret123")

			err := kr.Delete(tt.key)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				_, err := kr.Get(tt.key)
				require.Error(t, err)
			}
		})
	}
}
