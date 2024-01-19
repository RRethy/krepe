package threewaymerge

// twoWayMerge returns the result of performing a 2-way merge on local
// and upstream.
// If local and upstream are not the same type then local is returned.
// If local and upstream are the same type then the following rules apply:
//   - If local and upstream are maps then twoWayMergeMap algorithm is used.
//   - If local and upstream are slices then twoWayMergeSlice algorithm is used.
//   - If local and upstream are scalars then twoWayMergeScalar algorithm is used.
//
// local and upstream are not modified but the result may share memory with
// local and upstream.
func twoWayMerge(local, upstream any) any {
	switch localTyped := local.(type) {
	case map[string]any:
		upstreamMap, ok := upstream.(map[string]any)
		if !ok {
			return local
		}
		return twoWayMergeMap(localTyped, upstreamMap)
	case []any:
		upstreamArr, ok := upstream.([]any)
		if !ok {
			return local
		}
		return twoWayMergeSlice(localTyped, upstreamArr)
	default:
		return twoWayMergeScalar(local, upstream)
	}
}

// twoWayMergeMap returns the result of a recursive merge on each value in
// local using the corresponding value in upstream.
func twoWayMergeMap(local, upstream map[string]any) any {
	result := make(map[string]any)

	for k, v := range local {
		result[k] = v
	}

	for k, v := range upstream {
		if _, ok := result[k]; !ok {
			result[k] = v
		} else {
			result[k] = twoWayMerge(result[k], v)
		}
	}

	return result
}

// twoWayMergeSlice returns the result of performing a 2-way merge on local and
// upstream. The algorithm used depends on whether local and upstream are
// associative.
func twoWayMergeSlice(local, upstream []any) any {
	if isAssociativeSlice(local) && isAssociativeSlice(upstream) {
		return twoWayMergeSliceAssociative(local, upstream)
	}

	return twoWayMergeSliceNonAssociative(local, upstream)
}

// twoWayMergeSliceAssociative merges local and upstream using their associative
// key. An element present only in local will be returned as-is. An element
// present only in upstream will be returned as-is at the end of the slice. An
// element present in both local and upstream will be 2-way merged recursively.
func twoWayMergeSliceAssociative(local, upstream []any) any {
	key := getCommonAssociativeKey(getAssociativeKeys(local), getAssociativeKeys(upstream))
	if key == "" {
		return twoWayMergeSliceNonAssociative(local, upstream)
	}

	upstreamByKey := make(map[string]any, len(upstream))
	for _, elem := range upstream {
		elemMap, ok := elem.(map[string]any)
		if !ok {
			return twoWayMergeSliceNonAssociative(local, upstream)
		}

		keyVal, ok := elemMap[key].(string)
		if !ok {
			return twoWayMergeSliceNonAssociative(local, upstream)
		}

		upstreamByKey[keyVal] = elem
	}

	var result []any
	for _, elem := range local {
		elemMap, ok := elem.(map[string]any)
		if !ok {
			return twoWayMergeSliceNonAssociative(local, upstream)
		}

		keyVal, ok := elemMap[key].(string)
		if !ok {
			return twoWayMergeSliceNonAssociative(local, upstream)
		}

		if upstreamElem, ok := upstreamByKey[keyVal]; ok {
			result = append(result, twoWayMerge(elem, upstreamElem))
			delete(upstreamByKey, keyVal)
		} else {
			result = append(result, elem)
		}
	}

	for _, elem := range upstream {
		elemMap, ok := elem.(map[string]any)
		if !ok {
			return twoWayMergeSliceNonAssociative(local, upstream)
		}

		keyVal, ok := elemMap[key].(string)
		if !ok {
			return twoWayMergeSliceNonAssociative(local, upstream)
		}

		if _, ok := upstreamByKey[keyVal]; ok {
			result = append(result, elem)
		}
	}

	return result
}

// twoWayMergeSliceNonAssociative merges local and upstream using the 2-way
// merge scalar algorithm.
func twoWayMergeSliceNonAssociative(local, upstream []any) any {
	return twoWayMergeScalar(local, upstream)
}

// twoWayMergeScalar returns the upstream value.
func twoWayMergeScalar(local, upstream any) any {
	return upstream
}
