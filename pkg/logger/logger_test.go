package logger

import (
	"bytes"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestExplicitLogger_Log(t *testing.T) {
	const (
		repeats        = 3
		trivialMessage = "this is a trivial log message"
	)
	testCases := []struct {
		msg    string
		name   string
		format string
	}{
		{
			msg:    trivialMessage,
			name:   "Basic",
			format: time.RFC3339,
		},
		{
			msg:    "",
			name:   "EmptyMessage",
			format: time.RFC3339,
		},
		{
			msg:    trivialMessage,
			name:   "InvalidTimestamp",
			format: "invalid",
		},
		{
			msg:    trivialMessage,
			name:   "EmptyTimestamp",
			format: "",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			writer := new(bytes.Buffer)
			logger := NewExplicitLogger(writer, tc.msg, tc.name, tc.format)
			for i := 0; i < repeats; i++ {
				writer.Reset()
				ti, _, err := logger.Log()
				if err != nil {
					t.Fatalf("unexpected failure at message %d: %s", i, err.Error())
				}
				timestamp := ti.Format(tc.format)
				if len(timestamp) > 0 {
					timestamp += " "
				}
				want := fmt.Sprintf("%s[%s] %s\n", timestamp, tc.name, tc.msg)
				if want != writer.String() {
					t.Fatalf(`unexpected content: want "%s"; got "%s"`, want, writer.String())
				}
			}
		})
	}
}

func TestRandomLogger_Log(t *testing.T) {
	const repeats = 3
	testCases := []struct {
		size   int
		name   string
		format string
	}{
		{
			size:   256,
			name:   "Basic",
			format: time.RFC3339,
		},
		{
			size:   36,
			name:   "MinimalSize",
			format: time.RFC3339,
		},
		{
			size:   256,
			name:   "InvalidTimestamp",
			format: "invalid",
		},
		{
			size:   256,
			name:   "EmptyTimestamp",
			format: "",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			writer := new(bytes.Buffer)
			logger := NewRandomLogger(writer, tc.size, tc.name, tc.format)
			randStrings := make([]string, 0, repeats)
			for i := 0; i < repeats; i++ {
				writer.Reset()
				ti, _, err := logger.Log()
				if err != nil {
					t.Fatalf("unexpected failure at message %d: %s", i, err.Error())
				}
				timestamp := ti.Format(tc.format)
				if len(timestamp) > 0 {
					timestamp += " "
				}
				want := fmt.Sprintf("%s[%s]\n", timestamp, tc.name)
				compareRange := len(want) - 1
				if want[:compareRange] != writer.String()[:compareRange] {
					t.Fatalf(
						`unexpected content: want "%s..."; got "%s..."`,
						want[:compareRange], writer.String()[:compareRange],
					)
				}
				randStrings = append(randStrings, writer.String()[len(want):writer.Len()-1])
			}
			for l, r := 0, 1; r < repeats; l, r = r, r+1 {
				if len(randStrings[l]) > 0 {
					if randStrings[l] == randStrings[r] {
						t.Fatalf(
							`two consecutive logs have same content: left "%s" right "%s"`,
							randStrings[l], randStrings[r],
						)
					}
				}
			}
		})
	}
}

func BenchmarkRandomLogger_Log(b *testing.B) {
	const (
		size  = 1048576
		count = 500
	)
	writer := new(bytes.Buffer)
	writer.Grow(size)
	logger := NewRandomLogger(writer, size, "Benchmark", time.RFC3339)
	start := time.Now()
	b.ResetTimer()
	for i := 0; i < count; i++ {
		writer.Reset()
		b.StartTimer()
		if _, _, err := logger.Log(); err != nil {
			b.Fatalf("benchmark failed: %s", err.Error())
		}
	}
	b.StopTimer()
	log.Println(time.Now().Sub(start).String())
}
