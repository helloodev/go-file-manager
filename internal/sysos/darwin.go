package sysos

type DarwinProvider struct{}

func (d *DarwinProvider) GetOSName() string                         { return "macOS (Inoperative)"}
func (d *DarwinProvider) IsHidden(filename string) bool             { return false }
func (d *DarwinProvider) GetDiskUsage(path string) (string, error)  { return "N/A", nil }
                                                                    