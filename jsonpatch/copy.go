package jsonpatch

import (
	"fmt"
	"strconv"

	"golang.design/x/reflect"
)

type Copy struct {
	from []string
	path []string
}

func NewCopy(from, path string) (*Copy, error) {
	fromArr, err := pathToArray(from)
	if err != nil {
		return nil, fmt.Errorf("error parsing from: %s", err)
	}
	pathArr, err := pathToArray(path)
	if err != nil {
		return nil, fmt.Errorf("error parsing path: %s", err)
	}

	return &Copy{
		from: fromArr,
		path: pathArr,
	}, nil
}

func (c *Copy) Apply(obj any) (any, error) {
	return copy(obj, c.from, c.path)
}

func copy(obj any, from, path []string) (any, error) {
	fromObj, err := get(obj, from)
	if err != nil {
		return obj, err
	}

	return add(obj, path, reflect.DeepCopy(fromObj))
}

func get(obj any, path []string) (any, error) {
	if len(path) == 0 {
		return obj, nil
	}

	switch obj.(type) {
	case map[string]any:
		return getInMap(obj.(map[string]any), path[0], path[1:])
	case []any:
		return getInArray(obj.([]any), path[0], path[1:])
	default:
		return nil, fmt.Errorf("cannot get from non-map or non-array object with a json path")
	}
}

func getInMap(obj map[string]any, ptr string, path []string) (any, error) {
	val, ok := obj[ptr]
	if ok {
		return get(val, path)
	}

	return nil, fmt.Errorf("key not found")
}

func getInArray(obj []any, ptr string, path []string) (any, error) {
	idx, err := strconv.Atoi(ptr)
	if err != nil {
		return nil, fmt.Errorf("array index must be an integer")
	}

	if idx < 0 || idx >= len(obj) {
		return nil, fmt.Errorf("index out of range")
	}

	return get(obj[idx], path)
}
