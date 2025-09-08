package utils

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

// Format ke "%s jam %s menit"
func FormatIntervalToText(iv pgtype.Interval) string {
	totalSeconds := iv.Microseconds / 1_000_000
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60

	if hours > 0 && minutes > 0 {
		return fmt.Sprintf("%d jam %d menit", hours, minutes)
	} else if hours > 0 {
		return fmt.Sprintf("%d jam", hours)
	} else {
		return fmt.Sprintf("%d menit", minutes)
	}
}
