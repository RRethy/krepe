package jsonpatch

import (
	"fmt"
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

func (c *Test) Apply(obj any) (any, error) {
	return nil, nil
}
