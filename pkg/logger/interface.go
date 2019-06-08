package logger

import "time"

// Logger generates log message in various ways one at a time.
type Logger interface {
	// Log returns the time used to create the timestamp for the log and the size of the log, or an error if
	// any occurred.
	Log() (time.Time, int, error)
}
