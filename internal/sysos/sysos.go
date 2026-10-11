package sysos

import (
	"fmt"
	"runtime"
)

type OSAdapter interface {
	GetOSName()               string
	IsHidden(filename string) bool
	GetDiskUsage(path string) (string, error)
}

var CurrentOS OSAdapter

func InitOS() {
	switch runtime.GOOS {
	case "linux":
		CurrentOS = &LinuxProvider{}
		fmt.Println("🐧 [SYS::OS] Linux environment detected — POSIX/syscall provider initialized.")
	case "windows":
		CurrentOS = &WindowsProvider{}
		fmt.Println("🪟 [SYS::OS] Windows environment detected — compatibility mode active.")
	case "darwin":
		CurrentOS = &DarwinProvider{}
		fmt.Println("🍎 [SYS::OS] macOS environment detected — compatibility mode active.")
	default:
		CurrentOS = &LinuxProvider{}
		fmt.Println("⚙️ [SYS::OS] Unknown OS target — fallback POSIX provider initialized.")
	}
}