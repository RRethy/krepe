package jsonpatch

import (
	"fmt"
	"strconv"

	"github.com/RRethy/krepe/deepishcopy"
)

// Copy is a JSON patch operation that copies a value from one location to another within an object.
type Copy struct {
	from []string
	path []string
}

// NewCopy returns a Copy struct with the parsed from and path.
func NewCopy(from, path string) (*Copy, error) {
	fromArr, err := pathToArray(from)
	if err != nil {
		return nil, fmt.Errorf("parsing from: %s", err)
	}
	pathArr, err := pathToArray(path)
	if err != nil {
		return nil, fmt.Errorf("parsing path: %s", err)
	}

	return &Copy{
		from: fromArr,
		path: pathArr,
	}, nil
}

// Apply returns a new object that is the result of applying the Copy operation to obj.
func (c *Copy) Apply(obj any) (any, error) {
	obj, err := copy(obj, c.from, c.path)
	if err != nil {
		return nil, fmt.Errorf("copying from `%s` to `%s`: %s", c.from, c.path, err)
	}

	return obj, nil
}

func copy(obj any, from, path []string) (any, error) {
	fromObj, err := get(obj, from)
	if err != nil {
		return nil, err
	}

	return add(obj, path, deepishcopy.Copy(fromObj))
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
		return nil, fmt.Errorf("incompatible type for value %s: %T", obj, obj)
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
