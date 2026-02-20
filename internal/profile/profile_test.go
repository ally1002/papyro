package profile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ally1002/papyro/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		profile Profile
		wantErr bool
	}{
		{"all empty - should fail", Profile{}, true},
		{"only name - should fail", Profile{Name: "aly"}, true},
		{"only fromEmail - should fail", Profile{FromEmail: "aly@aly.com"}, true},
		{"only kindleEmail - should fail", Profile{KindleEmail: "aly@kindle.com"}, true},
		{"all fields - should pass", Profile{Name: "aly", FromEmail: "aly@aly.com", KindleEmail: "aly@kindle.com"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.profile.Validate()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func setup(t *testing.T) *Profiles {
	t.Helper()

	tempDir, err := os.MkdirTemp("", "papyro-test")
	require.NoError(t, err)

	t.Cleanup(func() { os.RemoveAll(tempDir) })

	cfg := &config.Config{
		Dir:      tempDir,
		FilePath: filepath.Join(tempDir, "profiles.json"),
	}
	cfg.CreateIfNotExists()

	return &Profiles{Profiles: []Profile{}, config: cfg}
}

func TestGet(t *testing.T) {
	ps := setup(t)
	err := ps.Add(&Profile{Name: "aly", FromEmail: "aly@aly.com", KindleEmail: "aly@kindle.com"})
	require.NoError(t, err)

	tests := []struct {
		name     string
		lookup   string
		wantErr  bool
		wantName string
	}{
		{"finds the profile by name - should pass", "aly", false, "aly"},
		{"does not find the profile by name - should fail", "alyen", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := ps.Get(tt.lookup)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, p)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantName, p.Name)
			}
		})
	}
}
