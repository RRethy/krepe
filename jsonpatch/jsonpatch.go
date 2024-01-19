// Package jsonpatch provides a JSON patch operation according to RFC 6902.
//
// It operates on Go data structures that would be produced from parsing JSON or YAML.
// It avoids the need to marshal and unmarshal JSON which most other JSON patch libraries require.
package jsonpatch

import (
	"fmt"
)

// JsonPatch is a JSON patch operation according to RFC 6902.
type JsonPatch interface {
	Apply(obj any) (any, error)
}

// NewJsonPatch returns a JsonPatch struct based on the given op.
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
