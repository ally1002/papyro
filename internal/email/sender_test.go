package email

import (
	"fmt"
	"os"
	"testing"

	"github.com/ally1002/papyro/internal/profile"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T) *Email {
	t.Helper()

	tempDir, _ := os.MkdirTemp("", "papyro-test")
	file, _ := os.CreateTemp(tempDir, "fileToSend_*.txt")

	profile := profile.Profile{Name: "aly", FromEmail: "aly@aly.com", KindleEmail: "aly@kindle.com"}
	password := []byte("12344321")
	filePath := file.Name()

	email, err := NewEmail(profile, password, filePath)
	if err != nil {
		t.Fatal(err)
	}

	return email
}

func TestNewEmail_ValidateFileExtension(t *testing.T) {
	profile := profile.Profile{Name: "aly", FromEmail: "aly@aly.com", KindleEmail: "aly@kindle.com"}
	password := []byte("12344321")

	tests := []struct {
		name    string
		ext     string
		wantErr bool
	}{
		{"accepts pdf  file - should pass", "pdf", false},
		{"accepts doc  file - should pass", "doc", false},
		{"accepts docx file - should pass", "docx", false},
		{"accepts txt  file - should pass", "txt", false},
		{"accepts rtf  file - should pass", "rtf", false},
		{"accepts htm  file - should pass", "htm", false},
		{"accepts html file - should pass", "html", false},
		{"accepts png  file - should pass", "png", false},
		{"accepts gif  file - should pass", "gif", false},
		{"accepts jpg  file - should pass", "jpg", false},
		{"accepts jpeg file - should pass", "jpeg", false},
		{"accepts bmp  file - should pass", "bmp", false},
		{"accepts epub file - should pass", "epub", false},
		{"does not accept nonlisted extensions - should fail", "webp", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir, _ := os.MkdirTemp("", "papyro-test")
			file, _ := os.CreateTemp(tempDir, fmt.Sprintf("fileToSend_*.%s", tt.ext))
			t.Cleanup(func() { _ = os.RemoveAll(tempDir) })

			filePath := file.Name()

			_, err := NewEmail(profile, password, filePath)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
