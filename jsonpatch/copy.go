package jsonpatch

import (
	"fmt"
)

type Copy struct {
	from []string
	path []string
}

func NewCopy(from, path string) (*Copy, error) {
	fromArr, err := pathToJsonPtrs(from)
	if err != nil {
		return nil, fmt.Errorf("error parsing from: %s", err)
	}
	pathArr, err := pathToJsonPtrs(path)
	if err != nil {
		return nil, fmt.Errorf("error parsing path: %s", err)
	}

	return &Copy{
		from: fromArr,
		path: pathArr,
	}, nil
}

func (c *Copy) Apply(obj any) (any, error) {
	return nil, nil
}
