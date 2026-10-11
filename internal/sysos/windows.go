package sysos

type WindowsProvider struct{}

func (w *WindowsProvider) GetOSName() string                             { return "Windows (Inoperative)"}
func (w *WindowsProvider) IsHidden(filename string) bool                 { return false }
func (w *WindowsProvider) GetDiskUsage(path string) ( string, error )    { return "N/A", nil }

