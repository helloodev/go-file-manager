package bot

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	
	"go-file-manager/internal/config"
	"go-file-manager/internal/fm"
	
	tele "gopkg.in/telebot.v3"
)

const MaxTelegramFileSize = 50 * 1024 * 1024

type UserAction struct {
	Action string
	Path   string
}

type Bot struct {
	tele   *tele.Bot
	cfg    *config.Config
	states map[int64]UserAction
	mu     sync.Mutex
}

func New(cfg *config.Config) (*Bot, error) {
	pref := tele.Settings{
		Token:  cfg.BotToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return nil, err
	}

	app := &Bot{
		tele: b,
		cfg:  cfg,
	}

	app.registerHandlers()
	return app, nil
}

func (b *Bot) setState(userID int64, action string, path string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.states == nil {
		b.states = make(map[int64]UserAction)
	}
	b.states[userID] = UserAction{Action: action, Path: path}
}

func (b *Bot) getState(userID int64) (UserAction, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.states == nil {
		return UserAction{}, false
	}
	act, ok := b.states[userID]
	return act, ok
}

func (b *Bot) clearState(userID int64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.states != nil {
		delete(b.states, userID)
	}
}

func (b *Bot) handleIncomingFile(c tele.Context, file *tele.File, fileName string) error {
	userID := c.Sender().ID
	state, ok := b.getState(userID)
	if !ok || state.Action != "upload" {
		return nil
	}

	destPath := filepath.Join(b.cfg.RootPath, state.Path, fileName)

	err := b.tele.Download(file, destPath)
	if err != nil {
		return c.Send(fmt.Sprintf("❌ [SYS_ERR]: Failed to save file: %v", err))
	}

	b.clearState(userID)
	_ = c.Send(fmt.Sprintf("✅ File uploaded successfully: `%s`", fileName), tele.ModeMarkdown)
	return fm.RenderFolder(c, b.cfg.RootPath, state.Path)
}

func (b *Bot) registerHandlers() {

	b.tele.Use(func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if c.Sender() == nil || c.Sender().ID != int64(b.cfg.OwnerID) {
				return nil
			}
			return next(c)
		}
	})
	
	b.tele.Handle("/start", func(c tele.Context) error {
		return fm.RenderFolder(c, b.cfg.RootPath, "")
	})

	b.tele.Handle(tele.OnDocument, func(c tele.Context) error {
		doc := c.Message().Document
		return b.handleIncomingFile(c, &doc.File, doc.FileName)
	})

	b.tele.Handle(tele.OnPhoto, func(c tele.Context) error {
		photo := c.Message().Photo
		fileName := fmt.Sprintf("photo_%d.jpg", time.Now().Unix())
		return b.handleIncomingFile(c, &photo.File, fileName)
	})

	b.tele.Handle(tele.OnVideo, func(c tele.Context) error {
		video := c.Message().Video
		fileName := video.FileName
		if fileName == "" {
			fileName = fmt.Sprintf("video_%d.mp4", time.Now().Unix())
		}
		return b.handleIncomingFile(c, &video.File, fileName)
	})

	b.tele.Handle(tele.OnAudio, func(c tele.Context) error {
		audio := c.Message().Audio
		fileName := audio.FileName
		if fileName == "" {
			fileName = fmt.Sprintf("audio_%d.mp3", time.Now().Unix())
		}
		return b.handleIncomingFile(c, &audio.File, fileName)
	})

	b.tele.Handle(tele.OnText, func(c tele.Context) error {
		userID := c.Sender().ID
		state, ok := b.getState(userID)
		if !ok {
			return nil
		}

		inputName := strings.TrimSpace(c.Text())
		if inputName == "" {
			return c.Send("❌ [SYS_ERR]: Name cannot be empty.")
		}

		b.clearState(userID)

		switch state.Action {
		case "mkdir":
			newFolderPath := filepath.Join(b.cfg.RootPath, state.Path, inputName)
			err := os.MkdirAll(newFolderPath, 0755)
			if err != nil {
				return c.Send(fmt.Sprintf("❌ [SYS_ERR]: Failed to create directory: %v", err))
			}
			_ = c.Send(fmt.Sprintf("✅ Directory created: `%s`", inputName), tele.ModeMarkdown)
			return fm.RenderFolder(c, b.cfg.RootPath, state.Path)

		case "rename":
			oldFullPath := filepath.Join(b.cfg.RootPath, state.Path)
			parentDir := filepath.Dir(state.Path)
			if parentDir == "." {
				parentDir = ""
			}
			newFullPath := filepath.Join(b.cfg.RootPath, parentDir, inputName)
			err := os.Rename(oldFullPath, newFullPath)
			if err != nil {
				return c.Send(fmt.Sprintf("❌ [SYS_ERR]: Failed to rename: %v", err))
			}
			_ = c.Send(fmt.Sprintf("✅ Renamed to: `%s`", inputName), tele.ModeMarkdown)
			return fm.RenderFolder(c, b.cfg.RootPath, parentDir)
		}
		return nil
	})


	b.tele.Handle(tele.OnCallback, func(c tele.Context) error {
		if c.Callback() == nil {
			return nil
		}

		data := strings.TrimPrefix(c.Callback().Data, "\f")
		parts := strings.SplitN(data, "|", 2)
		if len(parts) < 2 {
			return c.Respond()
		}

		action, pathID := parts[0], parts[1]
		path := fm.Store.GetPath(pathID)

		switch action {
		case "dir":
			b.clearState(c.Sender().ID)
			return fm.RenderFolder(c, b.cfg.RootPath, path)

		case "file":
			b.clearState(c.Sender().ID)
			return fm.RenderFile(c, b.cfg.RootPath, path)

		case "upload":
			b.setState(c.Sender().ID, "upload", path)
			return c.Send(fmt.Sprintf("📤 *[UPLOAD]*\nSend any file, photo, or document to save inside `/%s`:", path), tele.ModeMarkdown)
		
		case "mkdir":
			b.setState(c.Sender().ID, "mkdir", path)
			return c.Send(fmt.Sprintf("➕ *[MKDIR]*\nEnter new folder name inside `/%s`:", path), tele.ModeMarkdown)

		case "rename":
			b.setState(c.Sender().ID, "rename", path)
			return c.Send(fmt.Sprintf("✏️ *[RENAME]*\nEnter new name for `%s`:", filepath.Base(path)), tele.ModeMarkdown)

		case "send":
			fullPath := filepath.Join(b.cfg.RootPath, path)
			info, err := os.Stat(fullPath)
			if err != nil {
				return c.Respond(&tele.CallbackResponse{
					Text: "❌ File not found on host storage!",
					ShowAlert: true,
				})
			}

			if info.Size() > MaxTelegramFileSize {
				return c.Respond(&tele.CallbackResponse{
					Text: fmt.Sprintf("⚠️ FILE TOO LARGE (%s)!\nTelegram API limits direct uploads to 50MB.", fm.FormatSize(info.Size())),
					ShowAlert: true,
				})
			}
			
			doc := &tele.Document{
				File: tele.FromDisk(fullPath),
				FileName: filepath.Base(fullPath),
			}
			return c.Send(doc)

		case "del":
			promptText, err := fm.GetDeletionPrompt(b.cfg.RootPath, path)
			if err != nil {
				return c.Respond(&tele.CallbackResponse{
					Text:      "❌ Target not found!",
					ShowAlert: true,
				})
			}

			selector := &tele.ReplyMarkup{}
			btnConfirm := selector.Data("✅ YES, DELETE", "del_confirm", pathID)

			parentDir := filepath.Dir(path)
			if parentDir == "." {
				parentDir = ""
			}
			btnCancel := selector.Data("❌ CANCEL", "dir", fm.Store.GetID(parentDir))

			selector.Inline(selector.Row(btnConfirm, btnCancel))
			return c.EditOrSend(promptText, selector, tele.ModeMarkdown)

		case "del_confirm":
			err := fm.DeleteTarget(b.cfg.RootPath, path)
			if err != nil {
				return c.Send(fmt.Sprintf("❌ [SYS_ERR]: Failed to delete: %v", err))
			}

			parentDir := filepath.Dir(path)
			if parentDir == "." {
				parentDir = ""
			}
			promptText := fmt.Sprint("🗑 Target successfully deleted: `/%s`", path)
			_ = c.Send(promptText, tele.ModeMarkdown)
			return fm.RenderFolder(c, b.cfg.RootPath, parentDir)
		}

		return c.Respond()
	})
}

func (b *Bot) Start() {
	fmt.Printf("⚡ [SYS::INIT] File Manager Bot successfully launched...\n")
	b.tele.Start()
}