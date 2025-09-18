package dashboard

import (
	"fmt"
	"time"
)

func formatSpeed(bytes int) string {
	if bytes > 1024*1024 {
		return fmt.Sprintf("%.1f MB/s", float64(bytes)/(1024*1024))
	}
	return fmt.Sprintf("%.1f KB/s", float64(bytes)/1024)
}

func TruncateName(name string, max int) string {
	if len(name) <= max {
		return name
	}
	if max <= 3 {
		return name[:max]
	}
	return name[:max-3] + "..."
}

func FormatPercent(p float64) string {
	return fmt.Sprintf("%.2f%%", p*100)
}

func formatETA(seconds int) string {
	if seconds < 0 || seconds == 8640000 {
		return "—"
	}
	h := seconds / 3600
	m := (seconds % 3600) / 60

	if h > 0 {
		return fmt.Sprintf("%dh %dm", h, m)
	}
	return fmt.Sprintf("%dm", m)
}

func FormatAddedOn(ts int) string {
	if ts == 0 {
		return "—"
	}
	t := time.Unix(int64(ts), 0)
	return t.Format("2006-01-02 15:04")
}

func FormatSize(size int) string {
	sz := float64(size)
	switch {
	case sz >= 1<<30:
		return fmt.Sprintf("%.2f GB", sz/(1<<30))
	case sz >= 1<<20:
		return fmt.Sprintf("%.2f MB", sz/(1<<20))
	case sz >= 1<<10:
		return fmt.Sprintf("%.2f KB", sz/(1<<10))
	default:
		return fmt.Sprintf("%d B", size)
	}
}

func FormatBoolPtr(b *bool) string {
	if b == nil {
		return "—"
	}
	if *b {
		return "true"
	}
	return "false"
}