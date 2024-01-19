package jsonpatch

import (
	"fmt"
	"strings"
)

func JSONPatch(obj map[string]any, op, path string, value any) (map[string]any, error) {
	_, err := pathToJsonPtrs(path)
	if err != nil {
		return nil, err // TODO: wrap error
	}

	switch op {
	case "add":
		// newObj, err := add(obj, ptrs, value)
		// if err != nil {
		// return nil, err // TODO: wrap error
		// }
		// return newObj, nil
	}

	return obj, nil
}

func pathToJsonPtrs(path string) ([]string, error) {
	if path[0] != '/' {
		return nil, fmt.Errorf("path must start with /")
	}

	ptrs := strings.Split(path, "/")
	if len(ptrs) == 0 {
		return nil, fmt.Errorf("empty path is unsupported")
	}

	for i, ptr := range ptrs {
		ptr = strings.ReplaceAll(ptr, "~1", "/")
		ptr = strings.ReplaceAll(ptr, "~0", "~")
		ptrs[i] = ptr
	}

	return ptrs[1:], nil
}
