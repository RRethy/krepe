package jsonpatch

import (
	"fmt"
	"strconv"
)

type Replace struct {
	path  []string
	value any
}

func NewReplace(path string, value any) (*Replace, error) {
	pathArr, err := pathToArray(path)
	if err != nil {
		return nil, fmt.Errorf("error parsing path: %s", err)
	}

	return &Replace{
		path:  pathArr,
		value: value,
	}, nil
}

func (r *Replace) Apply(obj any) (any, error) {
	return replace(obj, r.path, r.value)
}

func replace(obj any, path []string, value any) (any, error) {
	if len(path) == 0 {
		return value, nil
	}

	switch obj.(type) {
	case map[string]any:
		return replaceMap(obj.(map[string]any), path[0], path[1:], value)
	case []any:
		return replaceArray(obj.([]any), path[0], path[1:], value)
	default:
		return obj, fmt.Errorf("cannot replace to non-map or non-array object with a json path")
	}
}

func replaceMap(obj map[string]any, ptr string, path []string, value any) (map[string]any, error) {
	if _, ok := obj[ptr]; !ok {
		return obj, fmt.Errorf("failed to replace to map: key not found: %s", ptr)
	}

	if len(path) == 0 {
		obj[ptr] = value
		return obj, nil
	}

	newVal, err := replace(obj[ptr], path, value)
	if err != nil {
		return obj, err
	}

	obj[ptr] = newVal
	return obj, nil
}

func replaceArray(arr []any, ptr string, path []string, value any) (any, error) {
	var idx int
	var err error
	if idx, err = strconv.Atoi(ptr); err != nil {
		return arr, err
	} else if idx < 0 || idx >= len(arr) {
		return arr, fmt.Errorf("index out of bounds: %d", idx)
	}

	if len(path) == 0 {
		arr[idx] = value
		return arr, nil
	}

	newVal, err := replace(arr[idx], path, value)
	if err != nil {
		return arr, err
	}

	arr[idx] = newVal
	return arr, nil
}
