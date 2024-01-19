package jsonpatch

import (
	"fmt"
	"reflect"
)

// Test is a JSON patch operation that asserts a value at a location within an object.
type Test struct {
	path  []string
	value any
}

// NewTest returns a Test struct with the parsed path and value.
func NewTest(path string, value any) (*Test, error) {
	pathArr, err := pathToArray(path)
	if err != nil {
		return nil, fmt.Errorf("parsing path: %s", err)
	}

	return &Test{
		path:  pathArr,
		value: value,
	}, nil
}

// Apply returns a new object that is the result of applying the Test operation to obj.
func (t *Test) Apply(obj any) (any, error) {
	got, err := get(obj, t.path)
	if err != nil {
		return nil, fmt.Errorf("getting value at path `%s`: %s", t.path, err)
	}

	if !reflect.DeepEqual(got, t.value) {
		return nil, fmt.Errorf("asserting value `%s` at path `%s`: got `%s`", t.value, t.path, got)
	}

	return obj, nil
}
