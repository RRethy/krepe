package jsonpatch

import (
	"fmt"
	"strconv"
)

func remove(obj any, ptrs []string) (any, error) {
	if len(ptrs) == 0 {
		return nil, nil
	}

	if idx, err := strconv.Atoi(ptrs[0]); err == nil {
		return removeArray(obj, idx, ptrs[1:])
	} else {
		return removeMap(obj, ptrs[0], ptrs[1:])
	}
}

func removeMap(obj any, ptr string, ptrs []string) (any, error) {
	var m map[string]any
	switch obj.(type) {
	case map[string]any:
		m = obj.(map[string]any)
	case nil:
		return obj, nil
	default:
		return obj, fmt.Errorf("cannot remove from non-map object with a non-number json ptr")
	}

	if len(ptrs) == 0 {
		delete(m, ptr)
		return m, nil
	}

	newVal, err := remove(m[ptr], ptrs)
	if err != nil {
		return obj, err
	}

	m[ptr] = newVal
	return m, nil
}

func removeArray(obj any, idx int, ptrs []string) (any, error) {
	var arr []any
	switch obj.(type) {
	case []any:
		arr = obj.([]any)
	case nil:
		return obj, nil
	default:
		return obj, fmt.Errorf("cannot remove from non-array object with a number json ptr")
	}

	if len(arr) <= idx {
		return obj, fmt.Errorf("cannot remove from array: index out of bounds")
	}

	if len(ptrs) == 0 {
		return append(arr[:idx], arr[idx+1:]...), nil
	}

	newVal, err := remove(arr[idx], ptrs)
	if err != nil {
		return obj, err
	}

	arr[idx] = newVal
	return arr, nil
}
