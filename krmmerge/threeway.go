package krmmerge

import (
	"reflect"
)

func threeWayMergeAny(origin, local, upstream any) any {
	switch localTyped := origin.(type) {
	case map[string]any:
		upstreamMap, ok := upstream.(map[string]any)
		if !ok {
			return local
		}

		localMap, ok := local.(map[string]any)
		if !ok {
			return twoWayMergeMap(localTyped, upstreamMap)
		}

		return threeWayMergeMap(localTyped, localMap, upstreamMap)
	case []any:
		upstreamArr, ok := upstream.([]any)
		if !ok {
			return local
		}

		localArr, ok := local.([]any)
		if !ok {
			return twoWayMergeSlice(localTyped, upstreamArr)
		}

		return threeWayMergeSlice(localTyped, localArr, upstreamArr)
	default:
		if reflect.TypeOf(local) != reflect.TypeOf(upstream) {
			return local
		}

		if reflect.TypeOf(local) != reflect.TypeOf(origin) {
			return twoWayMergeScalar(local, upstream)
		}

		return threeWayMergeScalar(origin, local, upstream)
	}

	return nil
}

func threeWayMergeMap(origin, local, upstream map[string]any) map[string]any {
	keepKeys := make(map[string]struct{})
	for k := range local {
		keepKeys[k] = struct{}{}
	}
	for k := range upstream {
		keepKeys[k] = struct{}{}
	}

	res := make(map[string]any)
	for key := range keepKeys {
		originVal, originOk := origin[key]
		localVal, localOk := local[key]
		upstreamVal, upstreamOk := upstream[key]

		if originOk && !localOk && !upstreamOk {
		} else if !originOk && localOk && !upstreamOk {
			res[key] = localVal
		} else if !originOk && !localOk && upstreamOk {
			res[key] = upstreamVal
		} else if originOk && localOk && !upstreamOk {
			val := delta(localVal, originVal)
			if val != nil {
				res[key] = val
			}
		} else if originOk && !localOk && upstreamOk {
			res[key] = upstreamVal
		} else if !originOk && localOk && upstreamOk {
			res[key] = twoWayMerge(localVal, upstreamVal)
		} else if originOk && localOk && upstreamOk {
			res[key] = threeWayMergeAny(originVal, localVal, upstreamVal)
		}
	}

	return res
}

func threeWayMergeSlice(origin, local, upstream []any) []any {
	if isAssociativeSlice(origin) && isAssociativeSlice(local) && isAssociativeSlice(upstream) {
		return threeWayMergeSliceAssociative(origin, local, upstream)
	}

	return threeWayMergeSliceNonAssociative(origin, local, upstream)
}

func threeWayMergeSliceAssociative(origin, local, upstream []any) []any {
	// key := getCommonAssociativeKey(getAssociativeKeys(origin), getAssociativeKeys(local), getAssociativeKeys(upstream))

	return nil
}

func threeWayMergeSliceNonAssociative(origin, local, upstream []any) []any {
	return nil
}

func threeWayMergeScalar(origin, local, upstream any) any {
	return nil
}
