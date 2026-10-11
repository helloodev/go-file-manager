package main

import (
	"log"
    
	"go-file-manager/internal/bot"
	"go-file-manager/internal/config"
	"go-file-manager/internal/sysos"
)


func main() {
    sysos.InitOS()

    cfg, err := config.LoadConfig("config.json")
    if err != nil {
        log.Fatalf("❌ [CONFIG_ERR]: %v", err)
    }

    app, err := bot.New(cfg)
    if err != nil {
        log.Fatalf("❌ [BOT_INIT_ERR]: %v", err)
    }

    app.Start()
}