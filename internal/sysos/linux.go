package sysos

import (
	"fmt"
	"strings"
	"syscall"
)

type LinuxProvider struct{}

func (l *LinuxProvider) GetOSName() string {
	return "Linux (POSIX)"
}

func (l *LinuxProvider) IsHidden(filename string) bool {
	return strings.HasPrefix(filename, ".")
}

func (l *LinuxProvider) GetDiskUsage(path string) (string, error) {
	var stat syscall.Statfs_t
	err := syscall.Statfs(path, &stat)
	if err != nil {
		return "", err
	}
	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bfree * uint64(stat.Bsize)
	used := total - free

	return fmt.Sprintf("Jami: %.2f | Ishlatildi: %.2f GB | Bo'sh: %.2f GB",
		float64(total)/(1024*1024*1024),
		float64(used)/(1024*1024*1024),
		float64(free)/(1024*1024*1024)), nil
}