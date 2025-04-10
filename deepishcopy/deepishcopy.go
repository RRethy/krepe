// Package deepishcopy provides a deep copy for data structures that would be produced from parsing YAML.
//
// It is purpose-built for krepe, but can be used for other purposes.
// Originally, we were using https://github.com/golang-design/reflect since it implements deep copy.
// However, it has a deal-breaking bug that arises when parsing YAML.
// See https://github.com/golang-design/reflect/issues/2 for more details.
package deepishcopy

import (
	"fmt"
	"reflect"
)

// Copy returns a deep copy of src.
//
// src must of type T where T is a:
//   - slice of T
//   - scalar
//   - nil
//   - map[K][V] where K is a string and V is of type T
func Copy(src any) any {
	switch src := src.(type) {
	case map[string]any:
		return copyMap(src)
	case []any:
		return copySlice(src)
	case bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64,
		float32, float64, complex64, complex128, string, nil:
		return src
	default:
		if reflect.TypeOf(src).Kind() == reflect.Ptr {
			return copyPtr(src)
		} else if reflect.TypeOf(src).Kind() == reflect.Struct {
			return copyStruct(src)
		}

		panic(fmt.Sprintf("deepishcopy: unsupported type %T", src))
	}
}

func copyMap(src map[string]any) map[string]any {
	if src == nil {
		return nil
	}

	dst := make(map[string]any, len(src))
	for k, v := range src {
		dst[k] = Copy(v)
	}
	return dst
}

func copySlice(src []any) []any {
	if src == nil {
		return nil
	}

	dst := make([]any, len(src))
	for i, v := range src {
		dst[i] = Copy(v)
	}
	return dst
}

func copyPtr(src any) any {
	srcVal := reflect.ValueOf(src)
	if srcVal.IsNil() {
		return reflect.Zero(reflect.TypeOf(src)).Interface()
	}

	dst := reflect.New(srcVal.Elem().Type()).Interface()
	dstElem := reflect.ValueOf(dst).Elem()
	srcElem := srcVal.Elem()

	if srcElem.CanInterface() && dstElem.CanSet() {
		dstElem.Set(reflect.ValueOf(Copy(srcElem.Interface())))
	}
	return dst
}

func copyStruct(src any) any {
	dst := reflect.New(reflect.TypeOf(src)).Elem()
	srcVal := reflect.ValueOf(src)

	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		dstField := dst.Field(i)

		if srcField.CanInterface() && dstField.CanSet() {
			srcFieldInterface := srcField.Interface()

			dstField.Set(reflect.ValueOf(Copy(srcFieldInterface)))
		}
	}
	return dst.Interface()
}
