package logger

import (
	"fmt"
	"io"
	"time"
)

type explicitLogger struct {
	writer          io.Writer
	name            string
	msg             string
	timestampFormat string
}

// NewExplicitLogger creates a Logger that prints a explicitly defined message.
func NewExplicitLogger(writer io.Writer, msg, name string, timestampFormat string) Logger {
	return &explicitLogger{
		writer:          writer,
		name:            name,
		msg:             msg,
		timestampFormat: timestampFormat,
	}
}

func (eg *explicitLogger) Log() (time.Time, int, error) {
	t, prefix := getPrefix(eg.name, eg.timestampFormat)
	size, err := eg.writer.Write([]byte(fmt.Sprintf("%s%s\n", prefix, eg.msg)))
	return t, size, err
}
