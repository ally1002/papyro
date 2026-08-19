package email

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/ally1002/papyro/internal/profile"
	"github.com/briandowns/spinner"
	"github.com/wneessen/go-mail"
)

type Mailer interface {
	DialAndSend(m ...*mail.Msg) error
}

type Email struct {
	profile  profile.Profile
	password []byte
	filePath string
	client   Mailer
}

var permittedExts = []string{"pdf", "doc", "docx", "txt", "rtf", "htm", "html", "png", "gif", "jpg", "jpeg", "bmp", "epub"}

func NewEmail(profile profile.Profile, password []byte, filePath string) (*Email, error) {
	if err := validateFile(filePath); err != nil {
		return nil, err
	}

	c, err := connectClient(profile, password)
	if err != nil {
		return nil, err
	}

	return &Email{profile: profile, password: password, filePath: filePath, client: c}, nil
}

func (e *Email) Send() error {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Sending email to Kindle...\n"

	s.Start()
	defer s.Stop()

	msg := mail.NewMsg()
	if err := msg.From(e.profile.FromEmail); err != nil {
		return fmt.Errorf("failed to set FROM address: %s", err)
	}
	if err := msg.To(e.profile.KindleEmail); err != nil {
		return fmt.Errorf("failed to set TO address: %s", err)
	}

	msg.Subject("Email sent by Papyro")
	msg.SetBodyString(mail.TypeTextPlain, "Here is your desired content.")
	msg.AttachFile(e.filePath)

	if err := e.client.DialAndSend(msg); err != nil {
		return fmt.Errorf("failed to deliver mail: %s", err)
	}

	s.FinalMSG = "Email sent successfully.\n"

	return nil
}

func connectClient(profile profile.Profile, password []byte) (*mail.Client, error) {
	return mail.NewClient("smtp.gmail.com",
		mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover), mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithUsername(profile.FromEmail), mail.WithPassword(string(password)),
		mail.WithTimeout(5*time.Minute),
	)
}

func validateFile(filePath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", filePath)
		}
		return fmt.Errorf("failed to access file: %s", err)
	}

	if info.IsDir() {
		return fmt.Errorf("path is a directory, not a file; supported formats: %s", formattedExts())
	}

	ext := strings.TrimPrefix(filepath.Ext(filePath), ".")
	if !isPermitted(ext) {
		return fmt.Errorf("unsupported file extension; supported formats: %s", formattedExts())
	}

	var maximumFileSize int64 = 50 * 1024 * 1024
	if info.Size() > maximumFileSize {
		return fmt.Errorf("file exceeds maximum size of 50 MB")
	}

	return nil
}

func isPermitted(ext string) bool {
	ext = strings.ToLower(ext)
	return slices.Contains(permittedExts, ext)
}

func formattedExts() string {
	formatted := make([]string, len(permittedExts))
	for i, ext := range permittedExts {
		formatted[i] = strings.ToUpper(ext)
	}
	return strings.Join(formatted, ", ")
}
