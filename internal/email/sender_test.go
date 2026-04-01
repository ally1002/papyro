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

func TestNewEmail_ValidateFileSize(t *testing.T) {
	profile := profile.Profile{Name: "aly", FromEmail: "aly@aly.com", KindleEmail: "aly@kindle.com"}
	password := []byte("12344321")

	tempDir, err := os.MkdirTemp("", "papyro-test")
	require.NoError(t, err, "failed to create temp dir")
	t.Cleanup(func() { _ = os.RemoveAll(tempDir) })

	tests := []struct {
		name     string
		fileSize int64
		wantErr  bool
	}{
		{"when file is not larger than 200mb - should pass", 25 * 1024 * 1024, false},
		{"when file is exactly 200mb - should pass", 200 * 1024 * 1024, false},
		{"when file is larger than 200mb - should fail", 200*1024*1024 + 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.CreateTemp(tempDir, "*.pdf")
			require.NoError(t, err)

			err = os.Truncate(file.Name(), tt.fileSize)
			require.NoError(t, err, "failed to create/truncate file")

			_, err = NewEmail(profile, password, file.Name())
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNewEmail_ValidateFileIsDir(t *testing.T) {
	profile := profile.Profile{Name: "aly", FromEmail: "aly@aly.com", KindleEmail: "aly@kindle.com"}
	password := []byte("12344321")

	tempDir, err := os.MkdirTemp("", "papyro-test")
	require.NoError(t, err, "failed to create temp dir")

	file, err := os.CreateTemp(tempDir, "fileToSend_*.txt")
	require.NoError(t, err, "failed to create temp file")

	t.Cleanup(func() { _ = os.RemoveAll(tempDir) })

	tests := []struct {
		name     string
		filePath string
		wantErr  bool
	}{
		{"when filePath is not a directory - should pass", file.Name(), false},
		{"when filePath is a directory - should fail", tempDir, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err = NewEmail(profile, password, tt.filePath)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}

}

func TestNewEmail_ValidateFileExistence(t *testing.T) {
	profile := profile.Profile{Name: "aly", FromEmail: "aly@aly.com", KindleEmail: "aly@kindle.com"}
	password := []byte("12344321")

	tempDir, err := os.MkdirTemp("", "papyro-test")
	require.NoError(t, err, "failed to create temp dir")

	file, err := os.CreateTemp(tempDir, "fileToSend_*.txt")
	require.NoError(t, err, "failed to create temp file")

	t.Cleanup(func() { _ = os.RemoveAll(tempDir) })

	tests := []struct {
		name     string
		filePath string
		wantErr  bool
	}{
		{"when file exists - should pass", file.Name(), false},
		{"when file does not exist - should fail", "file/does/not/exist.txt", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err = NewEmail(profile, password, tt.filePath)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}

}

func TestNewEmail_ValidateFileExtension(t *testing.T) {
	profile := profile.Profile{Name: "aly", FromEmail: "aly@aly.com", KindleEmail: "aly@kindle.com"}
	password := []byte("12344321")

	tempDir, err := os.MkdirTemp("", "papyro-test")
	require.NoError(t, err, "failed to create temp dir")
	t.Cleanup(func() { _ = os.RemoveAll(tempDir) })

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
			file, _ := os.CreateTemp(tempDir, fmt.Sprintf("fileToSend_*.%s", tt.ext))
			require.NoError(t, err, "failed to create temp file")

			filePath := file.Name()

			_, err = NewEmail(profile, password, filePath)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
