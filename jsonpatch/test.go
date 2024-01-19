package jsonpatch

import (
	"fmt"
	"reflect"
)

type Test struct {
	path  []string
	value any
}

func NewTest(path string, value any) (*Test, error) {
	pathArr, err := pathToArray(path)
	if err != nil {
		return nil, fmt.Errorf("error parsing path: %s", err)
	}

	return &Test{
		path:  pathArr,
		value: value,
	}, nil
}

func (t *Test) Apply(obj any) (any, error) {
	got, err := get(obj, t.path)
	if err != nil {
		return nil, err
	}

	if !reflect.DeepEqual(got, t.value) {
		return nil, fmt.Errorf("test failed: expected %v, got %v", t.value, got)
	}

	return obj, nil
}
