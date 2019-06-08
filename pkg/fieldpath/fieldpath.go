package fieldpath

import "strings"

// FieldPath is used to tract the path of the fields when validating a structure.
type FieldPath interface {
	// Returns the stored fields in a human-readable string.
	String() string

	// Return a copy of the calling FieldPath with the new field added; does not affect the original FieldPath.
	Add(string) FieldPath
}

type fieldPathImpl struct {
	fields []string
}

// NewFieldPath creates a new FieldPath, optionally with some base fields.
func NewFieldPath(fields ...string) FieldPath {
	return &fieldPathImpl{fields: fields}
}

func (f fieldPathImpl) String() string {
	return strings.Join(f.fields, ".")
}

func (f fieldPathImpl) Add(field string) FieldPath {
	ret := &fieldPathImpl{
		fields: make([]string, len(f.fields)+1),
	}
	copy(ret.fields, f.fields)
	ret.fields[len(f.fields)] = field
	return ret
}
