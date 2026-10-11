package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	BotToken string `json:"bot_token"`
	RootPath string `json:"root_path"`
	OwnerID int64 `json:"admin_id"`
}

func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	

	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	if cfg.RootPath == "" {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			cfg.RootPath = homeDir
		} else {
			cfg.RootPath = "./"
		}
	}

	return &cfg, nil
}
