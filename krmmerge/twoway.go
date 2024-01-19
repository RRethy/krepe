package krmmerge

func twoWayMerge(local, upstream any) any {
	switch local.(type) {
	case map[string]any:
		upstreamMap, ok := upstream.(map[string]any)
		if !ok {
			return local
		}
		return twoWayMergeMap(local.(map[string]any), upstreamMap)
	case []any:
		upstreamArr, ok := upstream.([]any)
		if !ok {
			return local
		}
		return twoWayMergeSlice(local.([]any), upstreamArr)
	default:
		return twoWayMergeScalar(local, upstream)
	}
}

func twoWayMergeMap(local, upstream map[string]any) map[string]any {
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

func twoWayMergeSlice(local, upstream []any) []any {
	if isAssociativeSlice(local) && isAssociativeSlice(upstream) {
		return twoWayMergeSliceAssociative(local, upstream)
	} else {
		return twoWayMergeSliceNonAssociative(local, upstream)
	}
}

func twoWayMergeSliceAssociative(local, upstream []any) []any {
	key := getAssociativeKey(local)
	if key == "" {
		return twoWayMergeSliceNonAssociative(local, upstream)
	}

	// result := make([]any, 0, len(local))
	// for _, localElem := range local {

	return nil
}

func twoWayMergeSliceNonAssociative(local, upstream []any) []any {
	return nil
}

func twoWayMergeScalar(local, upstream any) any {
	return nil
}
