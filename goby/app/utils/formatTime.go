package utils

import "time"

func FormatTime(input time.Time) string {
	return input.Format("2006-01-02 15:04:05")
}
