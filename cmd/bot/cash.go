package main

// import (
// 	"fmt"
// 	"go-file-manager/internal/config"
// 	"go-file-manager/internal/fm"
// 	"go-file-manager/internal/sysos"
// 	"log"
// 	"os"
// 	"path/filepath"
// 	"strings"
// 	"time"

// 	tele "gopkg.in/telebot.v3"
// )

// func main() {
//     sysos.InitOS()

//     cfg, err := config.LoadConfig("config.json")
//     if err != nil {
//         log.Fatalf("❌ Config faylini o'qishda xatolik: %v", err)
//     }

//     pref := tele.Settings{
//         Token: cfg.BotToken,
//         Poller: &tele.LongPoller{Timeout: 10 * time.Second},
//     }

//     b, err := tele.NewBot(pref)
//     if err != nil {
//         log.Fatal(err)
//     }

//     b.Handle("/start", func(c tele.Context) error{
//         return fm.RenderFolder(c, cfg.RootPath, "")
//         })

//     b.Handle(tele.OnCallback, func(c tele.Context) error{
//         if c.Callback() == nil {
//             return nil
//         }

//         data := strings.TrimPrefix(c.Callback().Data, "\f")
//         parts := strings.SplitN(data, "|", 2)
//         if len(parts) < 2 {
//             return c.Respond()
//         }

//         action, pathID := parts[0], parts[1]

//         path := fm.Store.GetPath(pathID)

//         switch action {
//         case "dir":
//             return fm.RenderFolder(c, cfg.RootPath, path)
//         case "file":
//             return fm.RenderFile(c, cfg.RootPath, path)
//         case "send":
//             fullPath := filepath.Join(cfg.RootPath, path)
//             doc := &tele.Document{
//                 File: tele.FromDisk(fullPath),
//                 FileName: filepath.Base(fullPath),
//             }
//             return c.Send(doc)
//         case "del":
//             fullPath := filepath.Join(cfg.RootPath)
//             _ = os.Remove(fullPath)
//             parentDir := filepath.Dir(path)
//             if parentDir == "." {
//                 parentDir = ""
//             }
//             return fm.RenderFolder(c, cfg.RootPath, parentDir)
//         }
//         return c.Respond()
//     })

//     fmt.Printf("🚀 File Manager Bot [%s] rejimida ishga tushdi...\n", sysos.CurrentOS.GetOSName())
//     b.Start()
// }