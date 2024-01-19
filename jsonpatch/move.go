package jsonpatch

import (
	"fmt"
)

type Move struct {
	from []string
	path []string
}

func NewMove(from, path string) (*Move, error) {
	fromPath, err := pathToArray(from)
	if err != nil {
		return nil, fmt.Errorf("error parsing from: %s", err)
	}

	toPath, err := pathToArray(path)
	if err != nil {
		return nil, fmt.Errorf("error parsing path: %s", err)
	}

	return &Move{
		from: fromPath,
		path: toPath,
	}, nil
}

func (m *Move) Apply(obj any) (any, error) {
	return move(obj, m.from, m.path)
}

func move(obj any, from []string, path []string) (any, error) {
	obj, removed, err := remove(obj, from)
	if err != nil {
		return nil, err
	}

	newObj, err := add(obj, path, removed)
	if err != nil {
		add(obj, from, removed)
		return nil, err
	}
	return newObj, nil
}
