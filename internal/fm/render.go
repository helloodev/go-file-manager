package fm

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go-file-manager/internal/sysos"

	tele "gopkg.in/telebot.v3"
)

func sanitizeMarkdown(text string) string {
	r := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"`", "\\`",
		"[", "\\[",
	)
	return r.Replace(text)
}

func GetDeletionPrompt(rootDir, relPath string) (string, error) {
	fullPath := filepath.Join(rootDir, relPath)
	info, err := os.Stat(fullPath)
	if err != nil {
		return "", err
	}

	if info.IsDir() {
		entries, err := os.ReadDir(fullPath)
		itemCount := 0
		if err == nil {
			itemCount = len(entries)
		}
		return fmt.Sprintf("⚠️ *[SYS::SECURITY_WARNING]*\nTarget: `/%s` (DIRECTORY)\n📁 Contents: `%d` items inside.\n\n_Deleting this folder will permanently purge all nested files and subdirectories. Proceed?_",
			sanitizeMarkdown(relPath), itemCount), nil
	}

	return fmt.Sprintf("⚠️ *[SYS::SECURITY_WARNING]*\nTarget: `/%s` (FILE)\n📊 Size: `%s`\n\n_Are you sure you want to permanently delete this file?_",
		sanitizeMarkdown(relPath), FormatSize(info.Size())), nil
}

func DeleteTarget(rootDir, relPath string) error {
	fullPath := filepath.Join(rootDir, relPath)
	return os.RemoveAll(fullPath)
}

func RenderFolder(c tele.Context, rootDir, relPath string) error {
	items, err := ReadDirectory(rootDir, relPath)
	if err != nil {
		return c.Send("❌ [SYS_ERR]: Failed to read directory.")
	}

	selector := &tele.ReplyMarkup{}
	var rows []tele.Row

	if relPath != "" {
		parent := filepath.Dir(relPath)
		if parent == "." {
			parent = ""
		}
		btnRename := selector.Data("✏️ RENAME", "rename", Store.GetID(relPath))
		btnUp := selector.Data("⬆️ UP ..", "dir", Store.GetID(parent))
		btnDelete := selector.Data("🗑 DELETE", "del", Store.GetID(relPath))

		rows = append(rows, selector.Row(btnRename, btnUp, btnDelete))
	}

	count := 0
	for _, item := range items {
		if item.IsHidden {
			continue
		}
		if count >= 40 {
			break
		}

		if item.IsDir {
			btn := selector.Data("📁 "+item.Name+"/", "dir", Store.GetID(item.Path))
			rows = append(rows, selector.Row(btn))
		} else {
			label := fmt.Sprintf("📄 %s (%s)", item.Name, FormatSize(item.Size))
			btn := selector.Data(label, "file", Store.GetID(item.Path))
			rows = append(rows, selector.Row(btn))
		}
		count++
	}

	parentPath := filepath.Dir(relPath)
	if parentPath == "." {
		parentPath = ""
	}

	btnUpload := selector.Data("📤 UPLOAD", "upload", Store.GetID(relPath))
	// btnBack := selector.Data("⬅️ BACK", "dir", Store.GetID(parentPath))
	btnMkdir := selector.Data("➕ MKDIR", "mkdir", Store.GetID(relPath))

	rows = append(rows, selector.Row(btnUpload, btnMkdir)) // btnBack qo'shilsa
	selector.Inline(rows...)

	diskInfo, _ := sysos.CurrentOS.GetDiskUsage(rootDir)
	displayPath := relPath
	if displayPath == "" {
		displayPath = "ROOT"
	}

	text := fmt.Sprintf("💻 *SYS::FM_SHELL*\n📂 *PATH:* `%s`\n📊 *ITEMS:* `%d`\n💾 *DISK_STAT:* `%s`",
		sanitizeMarkdown(displayPath),
		len(items),
		sanitizeMarkdown(diskInfo),
	)

	return c.EditOrSend(text, selector, tele.ModeMarkdown)
}

func RenderFile(c tele.Context, rootDir, relPath string) error {
	fullPath := filepath.Join(rootDir, relPath)
	info, err := os.Stat(fullPath)
	if err != nil {
		return c.Send("❌ [SYS_ERR]: Target file not found.")
	}

	selector := &tele.ReplyMarkup{}
	relID := Store.GetID(relPath)

	btnSend := selector.Data("⬇️ FETCH FILE", "send", relID)
	btnRename := selector.Data("✏️ RENAME", "rename", relID)
	btnDelete := selector.Data("🗑 DELETE", "del", relID)

	parentDir := filepath.Dir(relPath)
	if parentDir == "." {
		parentDir = ""
	}
	btnBack := selector.Data("⬅️ BACK", "dir", Store.GetID(parentDir))

	selector.Inline(
		selector.Row(btnSend),
		selector.Row(btnRename, btnDelete),
		selector.Row(btnBack),
	)

	text := fmt.Sprintf("📄 *FILE_INFO*\n🏷 *NAME:* `%s`\n📍 *PATH:* `%s`\n📊 *SIZE:* `%s`\n⚙️ *EXT:* `%s`",
		sanitizeMarkdown(info.Name()),
		sanitizeMarkdown(relPath),
		FormatSize(info.Size()),
		filepath.Ext(info.Name()),
	)

	return c.EditOrSend(text, selector, tele.ModeMarkdown)
}