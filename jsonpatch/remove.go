package jsonpatch

import (
	"fmt"
	"strconv"
)

type Remove struct {
	path []string
}

func NewRemove(path string) (*Remove, error) {
	pathArr, err := pathToArray(path)
	if err != nil {
		return nil, fmt.Errorf("error parsing path: %s", err)
	}

	return &Remove{
		path: pathArr,
	}, nil
}

func (r *Remove) Apply(obj any) (any, error) {
	obj, _, err := remove(obj, r.path)
	return obj, err
}

func remove(obj any, path []string) (any, any, error) {
	if len(path) == 0 {
		return nil, nil, nil
	}

	switch obj.(type) {
	case map[string]any:
		return removeInMap(obj.(map[string]any), path[0], path[1:])
	case []any:
		return removeInArray(obj.([]any), path[0], path[1:])
	default:
		return nil, nil, fmt.Errorf("cannot remove from non-map or non-array object with a json path")
	}
}

func removeInMap(obj map[string]any, ptr string, path []string) (any, any, error) {
	removed, ok := obj[ptr]
	if !ok {
		return nil, nil, fmt.Errorf("key not found")
	}

	if len(path) == 0 {
		delete(obj, ptr)
		return obj, removed, nil
	}

	newVal, removed, err := remove(obj[ptr], path)
	if err != nil {
		return nil, nil, err
	}

	obj[ptr] = newVal
	return obj, removed, nil
}

func removeInArray(arr []any, ptr string, path []string) (any, any, error) {
	var idx int
	var err error
	if idx, err = strconv.Atoi(ptr); err != nil {
		return arr, nil, err
	} else if idx < 0 || idx >= len(arr) {
		return nil, nil, fmt.Errorf("index out of bounds: %d", idx)
	}

	if len(path) == 0 {
		return append(arr[:idx], arr[idx+1:]...), arr[idx], nil
	}

	newVal, removed, err := remove(arr[idx], path)
	if err != nil {
		return nil, nil, err
	}

	arr[idx] = newVal
	return arr, removed, nil
}
