package email

import (
	"fmt"
	"os"
	"testing"
	"time"

	smtpmock "github.com/mocktools/go-smtp-mock/v2"
	"github.com/wneessen/go-mail"

	"github.com/ally1002/papyro/internal/profile"
	"github.com/stretchr/testify/require"
)

func TestValidateFileSize(t *testing.T) {
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
		{"under 50mb passes", 25 * 1024 * 1024, false},
		{"exactly 50mb passes", 50 * 1024 * 1024, false},
		{"fails when over 50mb", 50*1024*1024 + 1, true},
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

func TestValidateFileIsDir(t *testing.T) {
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
		{"file path passes", file.Name(), false},
		{"fails when path is directory", tempDir, true},
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

func TestValidateFileExistence(t *testing.T) {
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
		{"existing file passes", file.Name(), false},
		{"fails when file does not exist", "file/does/not/exist.txt", true},
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

func TestValidateFileExtension(t *testing.T) {
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
		{"pdf passes", "pdf", false},
		{"doc passes", "doc", false},
		{"docx passes", "docx", false},
		{"txt passes", "txt", false},
		{"rtf passes", "rtf", false},
		{"htm passes", "htm", false},
		{"html passes", "html", false},
		{"png passes", "png", false},
		{"gif passes", "gif", false},
		{"jpg passes", "jpg", false},
		{"jpeg passes", "jpeg", false},
		{"bmp passes", "bmp", false},
		{"epub passes", "epub", false},
		{"fails for unsupported extension", "webp", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.CreateTemp(tempDir, fmt.Sprintf("fileToSend_*.%s", tt.ext))
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

func TestSend(t *testing.T) {
	profile := profile.Profile{Name: "aly", FromEmail: "aly@aly.com", KindleEmail: "aly@kindle.com"}
	password := []byte("12344321")

	tempDir, err := os.MkdirTemp("", "papyro-test")
	require.NoError(t, err, "failed to create temp dir")

	file, err := os.CreateTemp(tempDir, "fileToSend_*.txt")
	require.NoError(t, err, "failed to create temp file")

	_, err = file.Write([]byte("Hello world"))
	require.NoError(t, err)
	err = file.Close()
	require.NoError(t, err)

	t.Cleanup(func() { _ = os.RemoveAll(tempDir) })

	server := smtpmock.New(smtpmock.ConfigurationAttr{
		LogToStdout:              false,
		LogServerActivity:        false,
		MultipleMessageReceiving: true,
	})

	err = server.Start()
	require.NoError(t, err, "smtp server did not start correctly")

	testClient, err := mail.NewClient("127.0.0.1",
		mail.WithPort(server.PortNumber()),
		mail.WithTLSPortPolicy(mail.NoTLS),
		mail.WithHELO("localhost"),
	)
	require.NoError(t, err, "mail client not created")

	email := &Email{profile: profile, password: password, filePath: file.Name(), client: testClient}

	err = email.Send()
	require.NoError(t, err)

	err = testClient.Close()
	require.NoError(t, err)

	_, err = server.WaitForMessages(1, 1*time.Second)
	require.NoError(t, err, "timeout waiting for messages")

	messages := server.Messages()
	require.NotEmpty(t, messages, "No messages received by mock server")

	msg := messages[0]
	require.True(t, msg.Msg(), "message should be received")
	require.True(t, msg.Data(), "DATA command should be successful")
	require.True(t, msg.Mailfrom(), "MAIL FROM should be successful")
	require.True(t, msg.Rcptto(), "RCPT TO should be successful")

	msgContent := msg.MsgRequest()
	require.Contains(t, msgContent, "Here is your desired content.")
	require.Contains(t, msgContent, "Email sent by Papyro")
	require.Contains(t, msgContent, "aly@aly.com")
	require.Contains(t, msgContent, "aly@kindle.com")
}
