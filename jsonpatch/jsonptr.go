package jsonpatch

import (
	"fmt"
	"strings"
)

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
