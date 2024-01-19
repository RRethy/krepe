package jsonpatch

import (
	"fmt"
	"strings"
)

func pathToArray(path string) ([]string, error) {
	if path == "" {
		return nil, nil
	}

	var arr []string
	var cur []rune
	for _, c := range path {
		if c == '/' {
			arr = append(arr, string(cur))
			cur = nil
		} else {
			cur = append(cur, c)
		}
	}
	arr = append(arr, string(cur))

	if arr[0] != "" {
		return nil, fmt.Errorf("path must start with /")
	}

	for i, part := range arr {
		part = strings.ReplaceAll(part, "~1", "/")
		part = strings.ReplaceAll(part, "~0", "~")
		arr[i] = part
	}

	return arr[1:], nil
}
