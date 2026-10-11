package fm

import (
	"fmt"
	"os"
	"path/filepath"

	"go-file-manager/internal/sysos"
)

type FileItem struct {
	Name      string
	IsDir     bool
	Size      int64
	Path      string
	IsHidden  bool
}

func ReadDirectory(rootDir, relPath string) ([]FileItem, error) {
	fullpath := filepath.Join(rootDir, relPath)
	entries, err := os.ReadDir(fullpath)
	if err != nil {
		return nil, err
	}

	var items []FileItem
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		items = append(items, FileItem{
			Name:      entry.Name(),
			IsDir:     entry.IsDir(),
			Size:      info.Size(),
			Path:      filepath.Join(relPath, entry.Name()),
			IsHidden:  sysos.CurrentOS.IsHidden(entry.Name()),
		})
	}
	return items, nil
}

func FormatSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	} else if size < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(size)/1024)
	}
	return fmt.Sprintf("%.1f MB", float64(size)/(1024*1024))
}