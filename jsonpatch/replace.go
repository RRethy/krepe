package jsonpatch

import (
	"fmt"
	"strconv"
)

func replace(obj any, ptrs []string, value any) (any, error) {
	if len(ptrs) == 0 {
		return value, nil
	}

	switch obj.(type) {
	case map[string]any:
		return replaceMap(obj.(map[string]any), ptrs[0], ptrs[1:], value)
	case []any:
		return replaceArray(obj.([]any), ptrs[0], ptrs[1:], value)
	default:
		return obj, fmt.Errorf("cannot replace to non-map or non-array object with a json ptr")
	}
}

func replaceMap(obj map[string]any, ptr string, ptrs []string, value any) (map[string]any, error) {
	if _, ok := obj[ptr]; !ok {
		return obj, fmt.Errorf("failed to replace to map: key not found: %s", ptr)
	}

	if len(ptrs) == 0 {
		obj[ptr] = value
		return obj, nil
	}

	newVal, err := replace(obj[ptr], ptrs, value)
	if err != nil {
		return obj, err
	}

	obj[ptr] = newVal
	return obj, nil
}

func replaceArray(arr []any, ptr string, ptrs []string, value any) (any, error) {
	var idx int
	var err error
	if idx, err = strconv.Atoi(ptr); err != nil {
		return arr, err
	} else if idx < 0 || idx >= len(arr) {
		return arr, fmt.Errorf("index out of bounds: %d", idx)
	}

	if len(ptrs) == 0 {
		arr[idx] = value
		return arr, nil
	}

	newVal, err := replace(arr[idx], ptrs, value)
	if err != nil {
		return arr, err
	}

	arr[idx] = newVal
	return arr, nil
}
