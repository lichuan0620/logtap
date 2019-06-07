package logger

import (
	"encoding/hex"
	"io"
	"math/rand"
	"sync"
	"time"
)

type randomLogger struct {
	output          io.Writer
	name            string
	timestampFormat string
	logBuffer       []byte
	hexBuffer       []byte
	newLinePos      int
	mutex           sync.Mutex
}

// NewRandomLogger creates a Logger that prints random strings no smaller than the minimal size.
func NewRandomLogger(writer io.Writer, size int, name string, timestampFormat string) Logger {
	ret := &randomLogger{
		output:          writer,
		name:            name,
		timestampFormat: timestampFormat,
	}
	_, prefix := getPrefix(name, timestampFormat)
	prefixSize := len(prefix)
	if prefixSize >= size {
		size = prefixSize + 1
	}
	ret.newLinePos = size - 1
	ret.logBuffer = make([]byte, size)
	ret.hexBuffer = make([]byte, size/2)
	ret.mutex.Lock()
	go ret.refresh()
	return ret
}

func (rg *randomLogger) Log() (time.Time, int, error) {
	t, prefix := getPrefix(rg.name, rg.timestampFormat)
	size, err := rg.doLog(prefix)
	rg.mutex.Lock()
	go rg.refresh()
	return t, size, err
}

func (rg *randomLogger) doLog(prefix string) (int, error) {
	rg.mutex.Lock()
	defer rg.mutex.Unlock()
	copy(rg.logBuffer, prefix)
	return rg.output.Write(rg.logBuffer)
}

func (rg *randomLogger) refresh() {
	defer rg.mutex.Unlock()
	rand.Read(rg.hexBuffer)
	hex.Encode(rg.logBuffer, rg.hexBuffer)
	rg.logBuffer[rg.newLinePos] = '\n'
}
