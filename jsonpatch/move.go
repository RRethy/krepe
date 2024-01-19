package jsonpatch

import (
	"fmt"
)

// Move is a JSON patch operation that moves a value from one location to another within an object.
type Move struct {
	from []string
	path []string
}

// NewMove returns a Move struct with the parsed from and path.
func NewMove(from, path string) (*Move, error) {
	fromPath, err := pathToArray(from)
	if err != nil {
		return nil, fmt.Errorf("parsing from: %s", err)
	}

	toPath, err := pathToArray(path)
	if err != nil {
		return nil, fmt.Errorf("parsing path: %s", err)
	}

	return &Move{
		from: fromPath,
		path: toPath,
	}, nil
}

// Apply returns a new object that is the result of applying the Move operation to obj.
func (m *Move) Apply(obj any) (any, error) {
	obj, err := move(obj, m.from, m.path)
	if err != nil {
		return nil, fmt.Errorf("moving value from `%s` to `%s`: %s", m.from, m.path, err)
	}

	return obj, nil
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
