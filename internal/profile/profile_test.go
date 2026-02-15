package profile

import (
	"testing"

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

func TestGet(t *testing.T) {
	tests := []struct {
		name    string
		profile Profile
	}{}
}
