package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateIfNotExists(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "papyro-test")
	require.NoError(t, err, "failed to create temp dir")

	t.Cleanup(func() {
		_ = os.RemoveAll(tempDir)
	})

	config := &Config{
		Dir:      tempDir,
		FilePath: filepath.Join(tempDir, "profiles.json"),
	}

	err = config.CreateIfNotExists()
	require.NoError(t, err, "failed to run CreateIfNotExists()")

	gotBytes, err := os.ReadFile(config.FilePath)
	require.NoError(t, err, "failed to read config file")

	var got map[string][]any
	err = json.Unmarshal(gotBytes, &got)
	require.NoError(t, err, "invalid JSON in config file")

	profiles, ok := got["profiles"]
	assert.True(t, ok, "config file missing 'profiles' key")
	assert.Len(t, profiles, 0, "profiles length should be zero")
}

func TestExists(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "papyro-test")
	require.NoError(t, err, "failed to create temp dir")

	t.Cleanup(func() {
		_ = os.RemoveAll(tempDir)
	})

	config := &Config{
		Dir:      tempDir,
		FilePath: filepath.Join(tempDir, "profiles.json"),
	}

	assert.False(t, config.Exists(), "should not be created")

	err = config.CreateIfNotExists()
	require.NoError(t, err, "failed to run CreateIfNotExists()")

	assert.True(t, config.Exists(), "Exists() should return true")
}
