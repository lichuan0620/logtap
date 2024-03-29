package v1alpha1

import (
	"fmt"
	"time"
)

const (
	// TaskPresetStandard produces a load of 256 B/log, 10 logs/s, and 2.5 KiB/s.
	TaskPresetStandard = "Standard"

	// TaskPresetLong produces a load of 20 MiB/log, 0.5 log/s, and 10 Mib/s.
	TaskPresetLong = "Long"

	// TaskPresetFrequent produces a load of 256 B/log, 50000 log/s, and 12 Mib/s.
	TaskPresetFrequent = "Frequent"

	// TaskPresetRoast produces a load of 1 MiB/log, 40 log/s, and 40 Mib/s.
	TaskPresetRoast = "Roast"
)

var (
	presets = map[string]*LogTaskSpec{
		TaskPresetStandard: {
			OutputKind:      OutputKindStdErr,
			TimestampFormat: time.RFC3339,
			ContentType:     ContentTypeRandom,
			MinSize:         256,
			Interval:        0.1,
		},
		TaskPresetLong: {
			OutputKind:      OutputKindStdErr,
			TimestampFormat: time.RFC3339,
			ContentType:     ContentTypeRandom,
			MinSize:         20971520,
			Interval:        2.,
		},
		TaskPresetFrequent: {
			OutputKind:      OutputKindStdErr,
			TimestampFormat: time.RFC3339,
			ContentType:     ContentTypeRandom,
			MinSize:         256,
			Interval:        0.00002,
		},
		TaskPresetRoast: {
			OutputKind:      OutputKindStdErr,
			TimestampFormat: time.RFC3339,
			ContentType:     ContentTypeRandom,
			MinSize:         1048576,
			Interval:        0.025,
		},
	}
)

// GetLogTaskSpecPreset returns a LogTaskSpec preset by its name or an error if no preset has that name.
func GetLogTaskSpecPreset(preset string) (*LogTaskSpec, error) {
	ret, exist := presets[preset]
	if !exist {
		return nil, fmt.Errorf("%s is not a valid preset", preset)
	}
	return ret.DeepCopy(), nil
}
