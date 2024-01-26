package merger

import (
	"fmt"
	"reflect"

	"github.com/RRethy/krepe/krepe/pkg/pkg"
)

// Merge returns the result of performing a 3-way merge on the origin, local, and upstream maps.
// This method does not mutate origin, local, and upstream, but the result might share data structures.
func ThreeWayMerge[T Mergeable](origin, local, upstream T) (T, error) {
	merged := threeWayMerge(any(origin), any(local), any(upstream))
	mergedTyped, ok := merged.(T)
	if !ok {
		return local, fmt.Errorf("TODO: internal error casting merged value of type %T to expected type %T", merged, local)
	}
	return mergedTyped, nil
}

// threeWayMerge returns the result of performing a 3-way merge on origin, local, and upstream.
// If any of origin, local, and upstream are not the same type then local is returned.
// If origin is a different type, then a 2-way merge is performed on local and upstream.
// If origin, local, and upstream are the same type then the following rules apply:
//   - If origin, local, and upstream are maps then threeWayMergeMap algorithm is used.
//   - If origin, local, and upstream are slices then threeWayMergeSlice algorithm is used.
//   - If origin, local, and upstream are scalars then threeWayMergeScalar algorithm is used.
//
// origin, local, and upstream are not modified but the result may share memory with origin, local, and upstream.
func threeWayMerge(origin, local, upstream any) any {
	if upstream == nil {
		if origin == nil {
			return local
		}
		return nil
	}

	if local == nil {
		if origin == nil {
			return upstream
		}
		return nil
	}

	if origin == nil {
		return twoWayMerge(local, upstream)
	}

	switch localTyped := local.(type) {
	case map[string]any:
		upstreamMap, ok := upstream.(map[string]any)
		if !ok {
			return local
		}

		originMap, ok := origin.(map[string]any)
		if !ok {
			return twoWayMergeMap(localTyped, upstreamMap)
		}

		return threeWayMergeMap(originMap, localTyped, upstreamMap)
	case []any:
		upstreamArr, ok := upstream.([]any)
		if !ok {
			return local
		}

		originArr, ok := origin.([]any)
		if !ok {
			return twoWayMergeSlice(localTyped, upstreamArr)
		}

		return threeWayMergeSlice(originArr, localTyped, upstreamArr)
	case *pkg.Pkg:
		upstreamPkg, ok := upstream.(*pkg.Pkg)
		if !ok {
			return local
		}

		originPkg, ok := origin.(*pkg.Pkg)
		if !ok {
			return twoWayMergePkg(localTyped, upstreamPkg)
		}

		return threeWayMergePkg(originPkg, localTyped, upstreamPkg)
	case *pkg.Krepe:
		upstreamKrepe, ok := upstream.(*pkg.Krepe)
		if !ok {
			return local
		}

		originKrepe, ok := origin.(*pkg.Krepe)
		if !ok {
			return twoWayMergeKrepe(localTyped, upstreamKrepe)
		}

		return threeWayMergeKrepe(originKrepe, localTyped, upstreamKrepe)
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

// threeWayMergeMap returns the result of a recursive 3-way merge on each key in local or upstream.
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
		} else if !originOk && localOk && upstreamOk {
			res[key] = twoWayMerge(localVal, upstreamVal)
		} else if originOk && localOk && upstreamOk {
			res[key] = threeWayMerge(originVal, localVal, upstreamVal)
		}
	}

	return res
}

// threeWayMergeSlice returns the result of performing a 3-way merge on origin, local, and upstream.
// The algorithm used depends on whether origin, local, and upstream are associative.
func threeWayMergeSlice(origin, local, upstream []any) []any {
	if isAssociativeSlice(origin) && isAssociativeSlice(local) && isAssociativeSlice(upstream) {
		return threeWayMergeSliceAssociative(origin, local, upstream)
	}

	return threeWayMergeSliceNonAssociative(origin, local, upstream)
}

// threeWayMergeSliceAssociative merges origin, local, and upstream elements by their associative key.
// An element only present in local will be returned as-is.
// An element only present in upstream will be returned as-is at the end of the slice.
// An element present in both local and upstream will be 2-way merged recursively.
// An element present in origin, local, and upstream will be 3-way merged recursively.
func threeWayMergeSliceAssociative(origin, local, upstream []any) []any {
	key := getCommonAssociativeKey(getAssociativeKeys(origin), getAssociativeKeys(local), getAssociativeKeys(upstream))
	if key == "" {
		return threeWayMergeSliceNonAssociative(origin, local, upstream)
	}

	type twoElemTuple struct {
		origin   any
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

		if originElem, ok := keysToRemove[keyVal]; ok {
			delete(keysToRemove, keyVal)
			keysToMerge[keyVal] = twoElemTuple{
				origin:   originElem,
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
			delete(keysToMerge, keyVal)
			result = append(result, threeWayMerge(keepTuple.origin, elem, keepTuple.upstream))
		} else if addElem, ok := keysToAdd[keyVal]; ok {
			delete(keysToAdd, keyVal)
			result = append(result, twoWayMerge(elem, addElem))
		} else {
			result = append(result, elem)
		}
	}

	for _, elem := range keysToAdd {
		result = append(result, elem)
	}

	return result
}

// threeWayMergeSliceNonAssociative returns the upstream value iff upstream has diverged from origin, otherwise local is returned.
func threeWayMergeSliceNonAssociative(origin, local, upstream []any) []any {
	if reflect.DeepEqual(origin, upstream) {
		return local
	}

	return upstream
}

// threeWayMergeScalar returns the upstream value iff upstream has diverged from origin, otherwise local is returned.
func threeWayMergeScalar(origin, local, upstream any) any {
	if reflect.DeepEqual(origin, upstream) {
		return local
	}

	return upstream
}

func threeWayMergePkg(origin, local, upstream *pkg.Pkg) *pkg.Pkg {
	// name := threeWayMergeScalar(origin.Name, local.Name, upstream.Name)
	// krepe := threeWayMerge(origin.Krepe, local.Krepe, upstream.Krepe)

	return nil
}

func threeWayMergeKrepe(origin, local, upstream *pkg.Krepe) *pkg.Krepe {
	return nil
}
