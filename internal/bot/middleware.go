package bot

import (
	"go-file-manager/internal/config"
)

func IsAdmin(cfg *config.Config, userID int64) bool {
	return userID == cfg.OwnerID
}