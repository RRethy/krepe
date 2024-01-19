package krmmerge

import (
	"reflect"
)

func delta(source, remove any) any {
	if reflect.TypeOf(source) != reflect.TypeOf(remove) {
		return source
	}

	switch source.(type) {
	case map[string]any:
		removeMap, ok := remove.(map[string]any)
		if !ok {
			return source
		}
		return deltaMap(source.(map[string]any), removeMap)
	case []any:
		removeArr, ok := remove.([]any)
		if !ok {
			return source
		}
		return deltaSlice(source.([]any), removeArr)
	default:
		return deltaScalar(source, remove)
	}
}

func deltaMap(source, remove map[string]any) any {
	res := make(map[string]any)

	for k, v := range source {
		if removeVal, ok := remove[k]; ok {
			newVal := delta(v, removeVal)
			if newVal != nil {
				res[k] = newVal
			}
		} else {
			res[k] = v
		}
	}

	if len(res) == 0 {
		return nil
	}
	return res
}

func deltaSlice(source, remove []any) any {
	if isAssociativeSlice(source) && isAssociativeSlice(remove) {
		return deltaSliceAssociative(source, remove)
	} else {
		return deltaSliceNonAssociative(source, remove)
	}
}

// deltaSliceNonAssociative returns nil iff source and remove are the same,
// otherwise source is returned.
func deltaSliceNonAssociative(source, remove []any) any {
	if len(source) != len(remove) {
		return source
	}

	for i, v := range source {
		if !reflect.DeepEqual(v, remove[i]) {
			return source
		}
	}

	return nil
}

// deltaSliceAssociative removes the elements from source that are present in remove.
// The elements are compared by the value of the associate key, if no key is
// found then a deep comparison is done.
// All elements in source and remove must be maps. Any duplicate keys will
// result in non associative behavior.
func deltaSliceAssociative(source, remove []any) any {
	key := getCommonAssociativeKey(getAssociativeKeys(source), getAssociativeKeys(remove))
	if key == "" {
		return deltaSliceNonAssociative(source, remove)
	}

	removeByKey := make(map[string]map[string]any)
	for _, elem := range remove {
		elemMap, ok := elem.(map[string]any)
		if !ok {
			return deltaSliceNonAssociative(source, remove)
		}

		key, ok := elemMap[key].(string)
		if !ok {
			return deltaSliceNonAssociative(source, remove)
		}

		removeByKey[key] = elemMap
	}

	var result []any
	for _, elem := range source {
		elemMap, ok := elem.(map[string]any)
		if !ok {
			return deltaSliceNonAssociative(source, remove)
		}

		key, ok := elemMap[key].(string)
		if !ok {
			return deltaSliceNonAssociative(source, remove)
		}

		if _, ok := removeByKey[key]; !ok {
			result = append(result, elem)
		}
	}

	return result
}

func deltaScalar(source, remove any) any {
	if reflect.DeepEqual(source, remove) {
		return nil
	}
	return source
}
