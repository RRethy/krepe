package jsonpatch

import (
	"fmt"
	"strconv"
)

type Add struct {
	path  []string
	value any
}

func NewAdd(path string, value any) (*Add, error) {
	pathArr, err := pathToArray(path)
	if err != nil {
		return nil, fmt.Errorf("error parsing path: %s", err)
	}

	return &Add{
		path:  pathArr,
		value: value,
	}, nil
}

func (a *Add) Apply(obj any) (any, error) {
	return add(obj, a.path, a.value)
}

func add(obj any, path []string, value any) (any, error) {
	if len(path) == 0 {
		return value, nil
	}

	switch obj.(type) {
	case map[string]any:
		return addMap(obj.(map[string]any), path[0], path[1:], value)
	case []any:
		return addArray(obj.([]any), path[0], path[1:], value)
	default:
		return obj, fmt.Errorf("cannot add to non-map or non-array object with a json path")
	}
}

func addMap(obj map[string]any, ptr string, path []string, value any) (map[string]any, error) {
	if len(path) == 0 {
		obj[ptr] = value
		return obj, nil
	}

	newVal, err := add(obj[ptr], path, value)
	if err != nil {
		return obj, err
	}

	obj[ptr] = newVal
	return obj, nil
}

func addArray(arr []any, ptr string, path []string, value any) (any, error) {
	var idx int
	var err error
	if ptr == "-" {
		idx = len(arr)
	} else if idx, err = strconv.Atoi(ptr); err != nil {
		return arr, err
	} else if idx < 0 || idx > len(arr) {
		return arr, fmt.Errorf("index out of bounds: %d", idx)
	}

	if len(path) == 0 {
		return append(arr[:idx], append([]any{value}, arr[idx:]...)...), nil
	}

	if idx == len(arr) {
		return arr, fmt.Errorf("index out of bounds: %d", idx)
	}
	newVal, err := add(arr[idx], path, value)
	if err != nil {
		return arr, err
	}

	arr[idx] = newVal
	return arr, nil
}
