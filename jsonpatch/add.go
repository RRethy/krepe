package jsonpatch

import (
	"fmt"
	"strconv"
)

func add(obj any, ptrs []string, value any) (any, error) {
	if len(ptrs) == 0 {
		return value, nil
	}

	switch obj.(type) {
	case map[string]any:
		return addMap(obj.(map[string]any), ptrs[0], ptrs[1:], value)
	case []any:
		return addArray(obj.([]any), ptrs[0], ptrs[1:], value)
	default:
		return obj, fmt.Errorf("cannot add to non-map or non-array object with a json ptr")
	}
}

func addMap(obj map[string]any, ptr string, ptrs []string, value any) (map[string]any, error) {
	if len(ptrs) == 0 {
		obj[ptr] = value
		return obj, nil
	}

	newVal, err := add(obj[ptr], ptrs, value)
	if err != nil {
		return obj, err
	}

	obj[ptr] = newVal
	return obj, nil
}

func addArray(arr []any, ptr string, ptrs []string, value any) (any, error) {
	var idx int
	var err error
	if ptr == "-" {
		idx = len(arr)
	} else if idx, err = strconv.Atoi(ptr); err != nil {
		return arr, err
	} else if idx < 0 || idx > len(arr) {
		return arr, fmt.Errorf("index out of bounds: %d", idx)
	}

	if len(ptrs) == 0 {
		return append(arr[:idx], append([]any{value}, arr[idx:]...)...), nil
	}

	if idx == len(arr) {
		return arr, fmt.Errorf("index out of bounds: %d", idx)
	}
	newVal, err := add(arr[idx], ptrs, value)
	if err != nil {
		return arr, err
	}

	arr[idx] = newVal
	return arr, nil
}
