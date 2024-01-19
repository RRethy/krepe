package jsonpatch

import (
	"fmt"
)

// Move is a JSON patch operation that moves a value from one location to another.
type Move struct {
	from []string
	path []string
}

// NewMove parses the from and path and creates a new Move struct.
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

// Apply applies the Move operation to the given object, and returns the modified object.
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
