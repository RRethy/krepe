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

func threeWayMergeSlice(origin, local, upstream []any) any {
	if isAssociativeSlice(origin) && isAssociativeSlice(local) && isAssociativeSlice(upstream) {
		return threeWayMergeSliceAssociative(origin, local, upstream)
	}

	return threeWayMergeSliceNonAssociative(origin, local, upstream)
}

func threeWayMergeSliceAssociative(origin, local, upstream []any) any {
	key := getCommonAssociativeKey(getAssociativeKeys(origin), getAssociativeKeys(local), getAssociativeKeys(upstream))
	if key == "" {
		return threeWayMergeSliceNonAssociative(origin, local, upstream)
	}

	type twoElemTuple struct {
		local    any
		upstream any
	}
	keysToRemove := make(map[string]any)
	keysToMerge := make(map[string]twoElemTuple)
	keysToAdd := make(map[string]any)

	for _, elem := range origin {
		elemMap, ok := elem.(map[string]any)
		if !ok {
			return threeWayMergeSliceNonAssociative(origin, local, upstream)
		}

		keyVal, ok := elemMap[key].(string)
		if !ok {
			return threeWayMergeSliceNonAssociative(origin, local, upstream)
		}

		keysToRemove[keyVal] = elem
	}

	for _, elem := range upstream {
		elemMap, ok := elem.(map[string]any)
		if !ok {
			return threeWayMergeSliceNonAssociative(origin, local, upstream)
		}

		keyVal, ok := elemMap[key].(string)
		if !ok {
			return threeWayMergeSliceNonAssociative(origin, local, upstream)
		}

		if localElem, ok := keysToRemove[keyVal]; ok {
			delete(keysToRemove, keyVal)
			keysToMerge[keyVal] = twoElemTuple{
				local:    localElem,
				upstream: elem,
			}
		} else {
			keysToAdd[keyVal] = elem
		}
	}

	var result []any
	for _, elem := range local {
		elemMap, ok := elem.(map[string]any)
		if !ok {
			return threeWayMergeSliceNonAssociative(origin, local, upstream)
		}

		keyVal, ok := elemMap[key].(string)
		if !ok {
			return threeWayMergeSliceNonAssociative(origin, local, upstream)
		}

		if _, ok := keysToRemove[keyVal]; ok {
			continue
		}

		if keepTuple, ok := keysToMerge[keyVal]; ok {
			result = append(result, threeWayMergeAny(keepTuple.local, elem, keepTuple.upstream))
		} else if addElem, ok := keysToAdd[keyVal]; ok {
			result = append(result, twoWayMerge(elem, addElem))
		} else {
			result = append(result, elem)
		}
	}

	return result
}

func threeWayMergeSliceNonAssociative(origin, local, upstream []any) any {
	return threeWayMergeScalar(origin, local, upstream)
}

func threeWayMergeScalar(origin, local, upstream any) any {
	return upstream
}
