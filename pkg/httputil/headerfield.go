package httputil

import (
	"fmt"
	"strings"
)

// HeaderField stores the key and values of a header field and provide functions to format the values.
type HeaderField interface {
	// Key returns the key of the header field
	Key() string

	// Value returns a string formatted from the values of the header field
	Value() string

	// String returns a string representation of the HeaderField
	String() string
}

type headerFieldImpl struct {
	key    string
	values []string
}

// NewHeaderField creates a new HeaderField.
func NewHeaderField(key string, values ...string) HeaderField {
	return &headerFieldImpl{
		key:    key,
		values: values,
	}
}

func (h headerFieldImpl) Key() string {
	return h.key
}

func (h headerFieldImpl) Value() string {
	return strings.Join(h.values, ",")
}

func (h headerFieldImpl) String() string {
	return fmt.Sprintf("%s: %s", h.key, h.Value())
}
