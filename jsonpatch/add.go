package jsonpatch

import (
	"fmt"
	"strconv"
)

func add(obj any, ptrs []string, value any) (any, error) {
	if len(ptrs) == 0 {
		return value, nil
	}

	if _, err := strconv.Atoi(ptrs[0]); err == nil || ptrs[0] == "-" {
		return addArray(obj, ptrs[0], ptrs[1:], value)
	} else {
		return addMap(obj, ptrs[0], ptrs[1:], value)
	}
}

func addMap(obj any, ptr string, ptrs []string, value any) (any, error) {
	var m map[string]any
	switch obj.(type) {
	case map[string]any:
		m = obj.(map[string]any)
	case nil:
		m = make(map[string]any)
	default:
		return obj, fmt.Errorf("cannot add to non-map object with a non-number json ptr")
	}

	newVal, err := add(m[ptr], ptrs, value)
	if err != nil {
		return obj, err
	}

	m[ptr] = newVal
	return m, nil
}

func addArray(obj any, ptr any, ptrs []string, value any) (any, error) {
	var arr []any
	switch obj.(type) {
	case []any:
		arr = obj.([]any)
	case nil:
	default:
		return obj, fmt.Errorf("cannot add to non-array object with a number json ptr")
	}

	var idx int
	if ptr == "-" {
		idx = len(arr)
	} else if n, err := strconv.Atoi(ptr.(string)); err == nil {
		idx = n
	} else {
		return obj, fmt.Errorf("unexpected error parsing array index: %s", ptr)
	}

	if len(arr) <= idx {
		arr = append(arr, make([]any, idx-len(arr)+1)...)
	}

	newVal, err := add(arr[idx], ptrs, value)
	if err != nil {
		return obj, err
	}

	arr[idx] = newVal
	return arr, nil
}
