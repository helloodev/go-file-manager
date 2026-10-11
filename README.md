# Go File Manager Bot

A lightweight Telegram-based file management system written in Go. It gives you a secure, interactive browser for a machine's filesystem through Telegram, allowing an authorized admin to upload, browse, rename, delete, and fetch files from a configured root directory.

This project is designed for self-hosted use cases where you want a simple, private file interface without exposing a full web dashboard.

## Features

- Browse directories and files from Telegram
- Upload documents, photos, videos, and audio into a chosen folder
- Create folders on the host machine
- Rename files or directories
- Delete files and folders safely with confirmation prompts
- Download or fetch files directly from Telegram
- Admin-only access using a configured Telegram user ID
- OS-aware behavior for hidden files and disk statistics
- Built in Go with minimal dependencies and a lightweight architecture

## Why this project exists

Many internal tools need a quick way to manage files on a remote machine without opening SSH or building a custom web app. This bot turns a local directory into a Telegram-native file manager, making it useful for:

- personal NAS or home-server management
- remote inspection of files and folders
- quick uploads from mobile devices
- private administrative file access

## Demo workflow

1. Start the bot with `/start`
2. Browse the configured root directory using inline Telegram buttons
3. Create folders, rename items, or open a file
4. Upload new content directly from your phone or desktop Telegram client
5. Fetch or delete data with confirmation prompts

## Supported usage

- File browsing through Telegram inline keyboard UI
- Upload to any nested directory under the configured root path
- Manage file names and directory structure in-place
- Download files back to Telegram when needed

## Project structure

```text
.
├── cmd/
│   └── bot/
│       └── main.go
├── internal/
│   ├── bot/
│   │   ├── bot.go
│   │   └── middleware.go
│   ├── config/
│   │   └── config.go
│   ├── fm/
│   │   ├── manager.go
│   │   ├── render.go
│   │   └── store.go
│   └── sysos/
│       ├── darwin.go
│       ├── linux.go
│       ├── sysos.go
│       └── windows.go
├── config.json
├── go.mod
└── README.md
```

## Prerequisites

- Go 1.27 or newer
- A Telegram bot token from BotFather
- Your Telegram user ID for admin access
- A writable root directory on the host machine

## Installation

Clone the repo:

```bash
git clone https://github.com/your-username/go-file-manager.git
cd go-file-manager
```

Install dependencies:

```bash
go mod download
```

## Configuration

Edit the `config.json` file:

```json
{
  "bot_token": "YOUR_TELEGRAM_BOT_TOKEN",
  "root_path": "/absolute/path/to/managed/files",
  "admin_id": 123456789
}
```

### Configuration fields

- `bot_token`: Telegram bot token received from BotFather
- `root_path`: The root directory the bot is allowed to manage
- `admin_id`: Telegram user ID allowed to interact with the bot

> Important: keep your bot token private and do not commit real credentials to a public repository.

## Running the bot

Start the bot with:

```bash
go run ./cmd/bot
```

The app will initialize the OS-specific helpers and start polling Telegram for updates.

## Usage

After launching the application:

1. Open your Telegram bot
2. Send `/start`
3. The bot will render the root folder as an interactive directory tree
4. Use inline buttons to:
   - open folders
   - upload files
   - create directories
   - rename entries
   - delete targets
   - fetch files

## Security notes

- Access is restricted to the configured admin ID only.
- The bot only operates inside the configured root path.
- Deletion prompts are shown before destructive actions.
- Avoid exposing your bot token or using a broad root path with sensitive data.

## Example scenarios

- Manage a local media folder from your mobile Telegram app
- Maintain a shared storage directory for personal automation
- Serve as a small private admin tool for a remote server

## Roadmap

Potential improvements for future versions:

- per-user authorization and multi-admin support
- authentication and access control layers
- file preview for images and text files
- archive extraction and compression utilities
- metrics and audit logs
- web UI companion for non-Telegram environments

## Contributing

Contributions are welcome. If you want to improve the project:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run validation checks
5. Open a pull request with a clear description

## Contact

For questions, suggestions, or feature requests, open an issue or reach out through the repository's discussion or contact channels.

---

Built for self-hosted file access through Telegram using Go.
