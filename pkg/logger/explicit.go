package logger

import (
	"fmt"
	"io"
	"time"
)

type explicitLogger struct {
	writer          io.Writer
	timestampFormat string
	logBuffer       []byte
}

// NewExplicitLogger creates a Logger that prints a explicitly defined message.
func NewExplicitLogger(writer io.Writer, msg, name string, timestampFormat string) Logger {
	_, prefix := getPrefix(name, timestampFormat)
	return &explicitLogger{
		writer:          writer,
		timestampFormat: timestampFormat,
		logBuffer:       []byte(fmt.Sprintf("%s%s\n", prefix, msg)),
	}
}

func (eg *explicitLogger) Log() (time.Time, int, error) {
	t, timestamp := getTimestamp(eg.timestampFormat)
	copy(eg.logBuffer, timestamp)
	size, err := eg.writer.Write(eg.logBuffer)
	return t, size, err
}
