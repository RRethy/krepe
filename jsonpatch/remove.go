package jsonpatch

import (
	"fmt"
	"strconv"
)

func remove(obj any, ptrs []string) (any, error) {
	if len(ptrs) == 0 {
		return nil, nil
	}

	switch obj.(type) {
	case map[string]any:
		return removeMap(obj.(map[string]any), ptrs[0], ptrs[1:])
	case []any:
		return removeArray(obj.([]any), ptrs[0], ptrs[1:])
	default:
		return obj, fmt.Errorf("cannot remove from non-map or non-array object with a json ptr")
	}
}

func removeMap(obj map[string]any, ptr string, ptrs []string) (any, error) {
	if _, ok := obj[ptr]; !ok {
		return obj, fmt.Errorf("key not found")
	}

	if len(ptrs) == 0 {
		delete(obj, ptr)
		return obj, nil
	}

	newVal, err := remove(obj[ptr], ptrs)
	if err != nil {
		return obj, err
	}

	obj[ptr] = newVal
	return obj, nil
}

func removeArray(arr []any, ptr string, ptrs []string) (any, error) {
	var idx int
	var err error
	if idx, err = strconv.Atoi(ptr); err != nil {
		return arr, err
	} else if idx < 0 || idx >= len(arr) {
		return arr, fmt.Errorf("index out of bounds: %d", idx)
	}

	if len(ptrs) == 0 {
		return append(arr[:idx], arr[idx+1:]...), nil
	}

	newVal, err := remove(arr[idx], ptrs)
	if err != nil {
		return arr, err
	}

	arr[idx] = newVal
	return arr, nil
}
