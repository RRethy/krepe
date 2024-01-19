package jsonpatch

import (
	"fmt"
	"strings"
)

type opType string

const (
	addOp     opType = "add"
	removeOp  opType = "remove"
	replaceOp opType = "replace"
	moveOp    opType = "move"
	copyOp    opType = "copy"
	testOp    opType = "test"
)

type JsonPatch struct {
	op    opType
	path  []string
	value any
}

func NewJsonPatch(op, path string, value any) (*JsonPatch, error) {
	ptrs, err := pathToJsonPtrs(path)
	if err != nil {
		return nil, fmt.Errorf("error parsing path: %s", err)
	}

	switch opType(op) {
	case addOp:
	case removeOp:
	case replaceOp:
	case moveOp:
	case copyOp:
	case testOp:
	default:
		return nil, fmt.Errorf("unknown op: %s", op)
	}

	return &JsonPatch{
		op:    opType(op),
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
	case removeOp:
		obj, err = remove(obj, jp.path)
		if err != nil {
			return obj, fmt.Errorf("failed applying 'remove': %s", err)
		}
		return obj, nil
	case replaceOp:
		// obj, err = replace(obj, jp.path, jp.value)
		// if err != nil {
		// return obj, fmt.Errorf("failed applying 'replace': %s", err)
		// }
		// return obj, nil
		panic("unimplemented")
	case moveOp:
		panic("unimplemented")
	case copyOp:
		panic("unimplemented")
	case testOp:
		panic("unimplemented")
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
