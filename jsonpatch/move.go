package jsonpatch

import (
	"fmt"
)

type Move struct {
	from []string
	path []string
}

func NewMove(from, path string) (*Move, error) {
	fromPtrs, err := pathToJsonPtrs(from)
	if err != nil {
		return nil, fmt.Errorf("error parsing from: %s", err)
	}

	pathPtrs, err := pathToJsonPtrs(path)
	if err != nil {
		return nil, fmt.Errorf("error parsing path: %s", err)
	}

	return &Move{
		from: fromPtrs,
		path: pathPtrs,
	}, nil
}

func (m *Move) Apply(obj any) (any, error) {
	return move(obj, m.from, m.path)
}

func move(obj any, from []string, ptrs []string) (any, error) {
	panic("not implemented")
}
