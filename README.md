# Papyro

A CLI tool to send documents to your Kindle via email. Built for speed and simplicity.

## Why Papyro?

Sending files to Kindle through Amazon's web interface is slow and cumbersome. Papyro gives you a simple command to send any supported file directly from your terminal:

```bash
papyro send my-kindle ~/Documents/book.pdf
```

No browser. No manual uploads. Just send.

## Installation

### Go

```bash
go install github.com/ally1002/papyro@latest
```

### Pre-built binaries

Download the latest release for your platform from the [releases page](https://github.com/ally1002/papyro/releases).

## Quick Start

1. **Add a profile**

   ```bash
   papyro profile add \
     --name my-kindle \
     --kindle-email your-kindle@kindle.com \
     --from-email your-email@gmail.com \
     --passwd your-app-password
   ```

   > **Note**: For Gmail, you'll need an [App Password](https://support.google.com/accounts/answer/185833). Your regular password won't work with SMTP.

2. **Send a file**

   ```bash
   papyro send my-kindle ~/Documents/book.pdf
   ```

That's it. The file will arrive on your Kindle shortly.

## Supported Formats

- Documents: PDF, DOC, DOCX, TXT, RTF, HTML
- Images: PNG, GIF, JPG, JPEG, BMP
- eBooks: EPUB

## File Size Limits

Papyro sends over SMTP, so it inherits your mail provider's limits. Attachments are base64-encoded on the wire, which inflates them by about 33% — a provider's "25 MB message" is roughly an **18 MB file on disk**.

| Path | Limit | Notes |
|------|-------|-------|
| Gmail SMTP (Papyro) | 25 MB encoded (≈ 18 MB file) | Enforced by Gmail, not Papyro |
| Amazon Send-to-Kindle email | ~50 MB | Applies to any mail provider |
| Send to Kindle web/app | 200 MB | Outside Papyro's scope |

The effective ceiling is whichever comes first — with Gmail, that's ~18 MB.

## Commands

| Command | Description |
|---------|-------------|
| `papyro profile add` | Create a new profile |
| `papyro profile list` | List all saved profiles |
| `papyro profile delete` | Remove a profile |
| `papyro send` | Send a file to Kindle |

## Configuration

Profiles are stored in `~/.config/papyro/profiles.json`. Passwords are stored securely in your system's keyring.

## License

[MIT](./LICENSE)
