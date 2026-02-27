package profile

import (
	"bytes"
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

func TestAdd(t *testing.T) {
	ps := setup(t)

	tests := []struct {
		name             string
		profileToBeAdded *Profile
		wantErr          bool
	}{
		{"add a valid profile - should pass", &Profile{Name: "aly", FromEmail: "aly@aly.com", KindleEmail: "aly@kindle.com"}, false},
		{"does not add an invalid profile - should fail", &Profile{Name: "", FromEmail: "", KindleEmail: ""}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ps.Add(tt.profileToBeAdded)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	ps := setup(t)
	err := ps.Add(&Profile{Name: "aly", FromEmail: "aly@aly.com", KindleEmail: "aly@kindle.com"})
	require.NoError(t, err)
	err = ps.Add(&Profile{Name: "alyen", FromEmail: "alyen@alyen.com", KindleEmail: "alyen@kindle.com"})
	require.NoError(t, err)

	tests := []struct {
		name               string
		profileToBeDeleted string
		wantErr            bool
	}{
		{"deletes the profile - should pass", "alyen", false},
		{"does not delete a profile that does not exist", "ghost", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ps.Delete(tt.profileToBeDeleted)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, ps.Profiles, 1)

				_, err := ps.Get(tt.profileToBeDeleted)
				require.Error(t, err)
			}
		})
	}
}

func TestList(t *testing.T) {
	ps := setup(t)

	err := ps.Add(&Profile{Name: "aly", FromEmail: "aly@aly.com", KindleEmail: "aly@kindle.com"})
	require.NoError(t, err)

	var buf bytes.Buffer

	err = ps.List(&buf)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "NAME")
	assert.Contains(t, output, "aly")
	assert.Contains(t, output, "aly@aly.com")
	assert.Contains(t, output, "aly@kindle.com")
}
