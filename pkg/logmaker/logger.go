package logmaker

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"sync"
	"time"
)

// logger is used internally by LogMaker to generate logs.
type logger interface {
	Log() (int, error)
}

// explicitLogger implements the logger interface; it simply repeats the given message.
type explicitLogger struct {
	writer          io.Writer
	timestampFormat string
	logBuffer       []byte
}

func newExplicitLogger(writer io.Writer, msg, name string, timestampFormat string) logger {
	return &explicitLogger{
		writer:          writer,
		timestampFormat: timestampFormat,
		logBuffer:       []byte(fmt.Sprintf("%s%s\n", getPrefix(name, timestampFormat), msg)),
	}
}

func (eg *explicitLogger) Log() (int, error) {
	copy(eg.logBuffer, getTimestamp(eg.timestampFormat))
	return eg.writer.Write(eg.logBuffer)
}

// randomLogger implements the randomLogger interface; it generates log messages that contain random characters
// and are no smaller than a minimal size.
type randomLogger struct {
	output          io.Writer
	name            string
	timestampFormat string
	logBuffer       []byte
	hexBuffer       []byte
	newLinePos      int
	mutex           sync.Mutex
}

func newRandomLogger(writer io.Writer, size int, name string, timestampFormat string) logger {
	ret := &randomLogger{
		output:          writer,
		name:            name,
		timestampFormat: timestampFormat,
	}
	prefix := getPrefix(name, timestampFormat)
	prefixSize := len(prefix)
	if prefixSize >= size {
		size = prefixSize + 1
	}
	ret.newLinePos = size - 1
	ret.logBuffer = make([]byte, size)
	ret.hexBuffer = make([]byte, size/2)
	ret.refresh()
	return ret
}

func (rg *randomLogger) Log() (int, error) {
	rg.mutex.Lock()
	defer rg.mutex.Unlock()
	copy(rg.logBuffer, getPrefix(rg.name, rg.timestampFormat))
	go rg.refresh()
	return rg.output.Write(rg.logBuffer)
}

// refresh should be called after generating a log message in a new goroutine. It prepares the next message.
func (rg *randomLogger) refresh() {
	rg.mutex.Lock()
	defer rg.mutex.Unlock()
	rand.Read(rg.hexBuffer)
	hex.Encode(rg.logBuffer, rg.hexBuffer)
	rg.logBuffer[rg.newLinePos] = '\n'
}

func getPrefix(name string, timestampFormat string) string {
	return fmt.Sprintf("%s [%s] ", getTimestamp(timestampFormat), name)
}

func getTimestamp(timestampFormat string) string {
	return time.Now().UTC().Format(timestampFormat)
}
