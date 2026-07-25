package cli

import (
	"time"

	"github.com/Blaze757/routeflux/internal/displaytime"
)

func formatLocalTimestamp(value time.Time) string {
	return displaytime.Format(value, time.RFC3339)
}

func formatLocalTimestampString(value string) string {
	return displaytime.FormatString(value, time.RFC3339)
}
