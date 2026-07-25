package tui

import (
	"time"

	"github.com/Blaze757/routeflux/internal/displaytime"
)

func formatLocalTimestamp(value time.Time) string {
	return displaytime.Format(value, "2006-01-02 15:04:05 MST")
}
