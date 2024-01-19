package jsonpatch

import (
	"fmt"
	"strings"
)

type opString string

const (
	addOp opString = "add"
)

type JsonPatch struct {
	op    opString
	path  []string
	value any
}

func NewJsonPatch(op, path string, value any) (*JsonPatch, error) {
	ptrs, err := pathToJsonPtrs(path)
	if err != nil {
		return nil, fmt.Errorf("error parsing path: %s", err)
	}

	switch opString(op) {
	case addOp:
	default:
		return nil, fmt.Errorf("unknown op: %s", op)
	}

	return &JsonPatch{
		op:    opString(op),
		path:  ptrs,
		value: value,
	}, nil
}

func (jp *JsonPatch) Apply(obj any) (any, error) {
	var err error
	switch jp.op {
	case addOp:
		obj, err = add(obj, jp.path, jp.value)
		if err != nil {
			return obj, fmt.Errorf("failed applying 'add': %s", err)
		}
		return obj, nil
	default:
		panic("unreachable")
	}
}

func pathToJsonPtrs(path string) ([]string, error) {
	if path == "" {
		return nil, nil
	}

	var ptrs []string
	var curptr []rune
	for _, c := range path {
		if c == '/' {
			ptrs = append(ptrs, string(curptr))
			curptr = nil
		} else {
			curptr = append(curptr, c)
		}
	}
	ptrs = append(ptrs, string(curptr))

	if ptrs[0] != "" {
		return nil, fmt.Errorf("path must start with /")
	}

	for i, ptr := range ptrs {
		ptr = strings.ReplaceAll(ptr, "~1", "/")
		ptr = strings.ReplaceAll(ptr, "~0", "~")
		ptrs[i] = ptr
	}

	return ptrs[1:], nil
}
