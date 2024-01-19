package jsonpatch

import (
	"fmt"
)

type JsonPatch interface {
	Apply(obj any) (any, error)
}

func NewJsonPatch(op, from, path string, value any) (JsonPatch, error) {
	switch op {
	case "add":
		return NewAdd(path, value)
	case "remove":
		return NewRemove(path)
	case "replace":
		return NewReplace(path, value)
	case "move":
		return NewMove(from, path)
	case "copy":
		return NewCopy(from, path)
	case "test":
		return NewTest(path, value)
	default:
		return nil, fmt.Errorf("unknown op: %s", op)
	}
}
