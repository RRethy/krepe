package jsonpatch

import (
	"fmt"
	"strconv"
)

// Replace is a JSON patch operation that replaces a value at a location within an object.
type Replace struct {
	path  []string
	value any
}

// NewReplace returns a Replace struct with the parsed path and value.
func NewReplace(path string, value any) (*Replace, error) {
	pathArr, err := pathToArray(path)
	if err != nil {
		return nil, fmt.Errorf("parsing path: %s", err)
	}

	return &Replace{
		path:  pathArr,
		value: value,
	}, nil
}

// Apply returns a new object that is the result of applying the Replace operation to obj.
func (r *Replace) Apply(obj any) (any, error) {
	obj, err := replace(obj, r.path, r.value)
	if err != nil {
		return nil, fmt.Errorf("replacing value at path `%s` with `%s`: %s", r.path, r.value, err)
	}

	return obj, nil
}

func replace(obj any, path []string, value any) (any, error) {
	if len(path) == 0 {
		return value, nil
	}

	switch obj := obj.(type) {
	case map[string]any:
		return replaceInMap(obj, path[0], path[1:], value)
	case []any:
		return replaceInArray(obj, path[0], path[1:], value)
	default:
		return nil, fmt.Errorf("incompatible type for value %s: %T", obj, obj)
	}
}

func replaceInMap(obj map[string]any, ptr string, path []string, value any) (map[string]any, error) {
	if _, ok := obj[ptr]; !ok {
		return nil, fmt.Errorf("key not found: %s", ptr)
	}

	if len(path) == 0 {
		obj[ptr] = value
		return obj, nil
	}

	newVal, err := replace(obj[ptr], path, value)
	if err != nil {
		return nil, err
	}

	obj[ptr] = newVal
	return obj, nil
}

func replaceInArray(arr []any, ptr string, path []string, value any) (any, error) {
	var idx int
	var err error
	if idx, err = strconv.Atoi(ptr); err != nil {
		return nil, err
	} else if idx < 0 || idx >= len(arr) {
		return nil, fmt.Errorf("index out of bounds: %d", idx)
	}

	if len(path) == 0 {
		arr[idx] = value
		return arr, nil
	}

	newVal, err := replace(arr[idx], path, value)
	if err != nil {
		return nil, err
	}

	arr[idx] = newVal
	return arr, nil
}
