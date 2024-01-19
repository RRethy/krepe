// Package deepishcopy provides a deep copy for data structures that would be produced from parsing YAML.
//
// It is purpose-built for krepe, but can be used for other purposes.
// Originally, we were using https://github.com/golang-design/reflect since it implements deep copy.
// However, it has a deal-breaking bug that arises when parsing YAML.
// See https://github.com/golang-design/reflect/issues/2 for more details.
package deepishcopy

import (
	"fmt"
)

// Copy returns a deep copy of src.
//
// src must of type T where T is a:
//   - slice of T
//   - scalar
//   - nil
//   - map[K][V] where K is a string and V is of type T
func Copy(src any) any {
	switch src.(type) {
	case map[string]any:
		return copyMap(src.(map[string]any))
	case []any:
		return copySlice(src.([]any))
	case bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64,
		float32, float64, complex64, complex128, string, nil:
		return src
	default:
		panic(fmt.Sprintf("deepishcopy: unsupported type %T", src))
	}
}

func copyMap(src map[string]any) map[string]any {
	dst := make(map[string]any, len(src))
	for k, v := range src {
		dst[k] = Copy(v)
	}
	return dst
}

func copySlice(src []any) []any {
	dst := make([]any, len(src))
	for i, v := range src {
		dst[i] = Copy(v)
	}
	return dst
}
