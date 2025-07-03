package pbf

import (
	"fmt"
	"time"
)

// Utility functions

func generateTransactionCode() string {
	return fmt.Sprintf("TRX%d", time.Now().Unix())
}

func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

func parseDateTime(dateTimeStr string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", dateTimeStr)
}
