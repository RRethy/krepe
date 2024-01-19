package krmmerge

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

func twoWayMergeSlice(local, upstream []any) any {
	if isAssociativeSlice(local) && isAssociativeSlice(upstream) {
		return twoWayMergeSliceAssociative(local, upstream)
	} else {
		return twoWayMergeSliceNonAssociative(local, upstream)
	}
}

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

func twoWayMergeSliceNonAssociative(local, upstream []any) any {
	return twoWayMergeScalar(local, upstream)
}

func twoWayMergeScalar(local, upstream any) any {
	return upstream
}
