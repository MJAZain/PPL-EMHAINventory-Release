package utils

import (
	"strconv"
	"time"
)

// GetRelativeTimeString returns a relative time description in Indonesian.
func GetRelativeTimeString(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		return "baru saja"
	case diff < time.Hour:
		minutes := int(diff.Minutes())
		return pluralize(minutes, "menit")
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		return pluralize(hours, "jam")
	case diff < 30*24*time.Hour:
		days := int(diff.Hours() / 24)
		return pluralize(days, "hari")
	default:
		months := int(diff.Hours() / 24 / 30)
		return pluralize(months, "bulan")
	}
}

func pluralize(count int, unit string) string {
	if count == 1 {
		return "1 " + unit + " yang lalu"
	}
	return strconv.Itoa(count) + " " + unit + " yang lalu"
}
