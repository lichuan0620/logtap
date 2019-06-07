package logger

import (
	"fmt"
	"time"
)

func getPrefix(name string, timestampFormat string) (time.Time, string) {
	t, timestamp := getTimestamp(timestampFormat)
	if len(timestamp) > 0 {
		return t, fmt.Sprintf("%s [%s] ", timestamp, name)
	}
	return t, fmt.Sprintf("[%s] ", name)
}

func getTimestamp(timestampFormat string) (time.Time, string) {
	now := time.Now().UTC()
	return now, now.Format(timestampFormat)
}
