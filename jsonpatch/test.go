package jsonpatch

import (
	"fmt"
	"reflect"
)

// Test is a JSON patch operation that tests a value at a location.
type Test struct {
	path  []string
	value any
}

// NewTest parses the path and creates a new Test struct.
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

// Apply applies the Test operation to the given object, and returns the modified object.
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
