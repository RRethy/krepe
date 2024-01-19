package threewaymerge

import (
	"reflect"
)

// delta returns the difference between source and remove.
// If source and remove are not the same type then source is returned.
// If source and remove are the same type then the following rules apply:
//   - If source and remove are maps then deltaMap algorithm is used.
//   - If source and remove are slices then deltaSlice algorithm is used.
//   - If source and remove are scalars then deltaScalar algorithm is used.
//
// source and remove are not modified but the result may share memory with
// source and remove.
func delta(source, remove any) any {
	if reflect.TypeOf(source) != reflect.TypeOf(remove) {
		return source
	}

	switch sourceTyped := source.(type) {
	case map[string]any:
		removeMap, ok := remove.(map[string]any)
		if !ok {
			return source
		}

		res := deltaMap(sourceTyped, removeMap)
		if res == nil {
			return nil
		}
		return res
	case []any:
		removeArr, ok := remove.([]any)
		if !ok {
			return source
		}

		res := deltaSlice(sourceTyped, removeArr)
		if res == nil {
			return nil
		}
		return res
	default:
		res := deltaScalar(source, remove)
		if res == nil {
			return nil
		}
		return res
	}
}

// deltaMap returns the difference between source map and remove map.
// A recursive delta is performed on each value in source using the
// corresponding value in remove.
// If the result is an empty map then nil is returned.
func deltaMap(source, remove map[string]any) map[string]any {
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

// deltaSlice returns the difference between source slice and remove slice.
// The algorithm used depends on whether or not the slices are associative.
func deltaSlice(source, remove []any) []any {
	if isAssociativeSlice(source) && isAssociativeSlice(remove) {
		return deltaSliceAssociative(source, remove)
	}

	return deltaSliceNonAssociative(source, remove)
}

// deltaSliceAssociative removes the elements from source that are present in
// remove. Comparison is done by the value of the common associative key. Both
// source and remove must be associative slices or the non-associative
// algorithm is used.
func deltaSliceAssociative(source, remove []any) []any {
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

// deltaSliceNonAssociative returns nil iff source and remove are the same,
// otherwise source is returned.
func deltaSliceNonAssociative(source, remove []any) []any {
	if reflect.DeepEqual(source, remove) {
		return nil
	}
	return source
}

// deltaScalar returns nil iff source and remove are the same,
// otherwise source is returned.
func deltaScalar(source, remove any) any {
	if reflect.DeepEqual(source, remove) {
		return nil
	}
	return source
}
