package merger

import (
	"reflect"
)

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
	if origin == nil {
		return twoWayMerge(local, upstream)
	}

	if upstream == nil || local == nil {
		return threeWayMergeScalar(origin, local, upstream)
	}

	localType := reflect.TypeOf(local)
	upstreamType := reflect.TypeOf(upstream)
	originType := reflect.TypeOf(origin)

	if localType != upstreamType || localType != originType {
		return threeWayMergeScalar(origin, local, upstream)
	}

	switch localTyped := local.(type) {
	case map[string]any:
		return threeWayMergeMap(origin.(map[string]any), localTyped, upstream.(map[string]any))
	case []any:
		return threeWayMergeSlice(origin.([]any), localTyped, upstream.([]any))
	default:
		if localType.Kind() == reflect.Ptr && localType.Elem().Kind() == reflect.Struct {
			return threeWayMergePtrStruct(origin, local, upstream)
		} else if localType.Kind() == reflect.Struct {
			return threeWayMergeStruct(origin, local, upstream)
		}

		return threeWayMergeScalar(origin, local, upstream)
	}
}

// threeWayMergeMap returns the result of a recursive 3-way merge on each key in local or upstream.
func threeWayMergeMap(origin, local, upstream map[string]any) any {
	keepKeys := make(map[string]struct{})
	for k := range local {
		keepKeys[k] = struct{}{}
	}
	for k := range upstream {
		keepKeys[k] = struct{}{}
	}

	res := make(map[string]any)
	for key := range keepKeys {
		originVal := origin[key]
		localVal := local[key]
		upstreamVal := upstream[key]

		newVal := threeWayMerge(originVal, localVal, upstreamVal)
		if newVal != nil {
			res[key] = newVal
		}
	}

	if len(res) == 0 {
		return nil
	}
	return res
}

// threeWayMergeSlice returns the result of performing a 3-way merge on origin, local, and upstream.
// The algorithm used depends on whether origin, local, and upstream are associative.
func threeWayMergeSlice(origin, local, upstream []any) any {
	if isUniformStructSlices(origin, local, upstream) {
		return threeWayMergeStructSlice(origin, local, upstream)
	}

	if isUniformPtrStructSlices(origin, local, upstream) {
		return threeWayMergePtrStructSlice(origin, local, upstream)
	}

	if isAssociativeSlices(origin, local, upstream) {
		return threeWayMergeSliceAssociative(origin, local, upstream)
	}

	return threeWayMergeSliceNonAssociative(origin, local, upstream)
}

// threeWayMergeSliceAssociative merges origin, local, and upstream elements by their associative key.
// An element only present in local will be returned as-is.
// An element only present in upstream will be returned as-is at the end of the slice.
// An element present in both local and upstream will be 2-way merged recursively.
// An element present in origin, local, and upstream will be 3-way merged recursively.
func threeWayMergeSliceAssociative(origin, local, upstream []any) any {
	key, ok := getCommonAssociativeKey(getAssociativeKeys(origin), getAssociativeKeys(local), getAssociativeKeys(upstream))
	if !ok {
		return threeWayMergeSliceNonAssociative(origin, local, upstream)
	}

	type groupedItems struct {
		origin   any
		local    any
		upstream any
	}
	items := map[string]*groupedItems{}
	for _, elem := range origin {
		items[elem.(map[string]any)[key].(string)] = &groupedItems{origin: elem}
	}
	for _, elem := range local {
		key := elem.(map[string]any)[key].(string)
		if group, ok := items[key]; ok {
			group.local = elem
		} else {
			items[key] = &groupedItems{local: elem}
		}
	}
	for _, elem := range upstream {
		key := elem.(map[string]any)[key].(string)
		if group, ok := items[key]; ok {
			group.upstream = elem
		} else {
			items[key] = &groupedItems{upstream: elem}
		}
	}

	var res []any
	used := map[string]struct{}{}
	for _, elem := range local {
		key := elem.(map[string]any)[key].(string)
		group := items[key]
		newVal := threeWayMerge(group.origin, group.local, group.upstream)
		if newVal != nil {
			res = append(res, newVal)
		}
		used[key] = struct{}{}
	}
	for _, elem := range upstream {
		key := elem.(map[string]any)[key].(string)
		if _, ok := used[key]; !ok {
			group := items[key]
			newVal := threeWayMerge(group.origin, group.local, group.upstream)
			if newVal != nil {
				res = append(res, newVal)
			}
		}
	}

	return res
}

// threeWayMergeSliceNonAssociative returns the upstream value iff upstream has diverged from origin, otherwise local is returned.
func threeWayMergeSliceNonAssociative(origin, local, upstream []any) any {
	return threeWayMergeScalar(origin, local, upstream)
}

// threeWayMergeScalar returns the upstream value iff upstream has diverged from origin, otherwise local is returned.
func threeWayMergeScalar(origin, local, upstream any) any {
	if reflect.DeepEqual(origin, upstream) {
		return local
	}

	return upstream
}

func threeWayMergePtrStruct(origin, local, upstream any) any {
	if local == nil || upstream == nil {
		return threeWayMergeScalar(origin, local, upstream)
	}
	if origin == nil {
		return twoWayMergePtrStruct(local, upstream)
	}

	originMap := ptrStructToMap(origin)
	localMap := ptrStructToMap(local)
	upstreamMap := ptrStructToMap(upstream)

	merged := threeWayMergeMap(originMap, localMap, upstreamMap)
	if merged == nil {
		return reflect.New(reflect.TypeOf(local).Elem()).Interface()
	}

	targetType := threeWayMergeScalar(
		reflect.TypeOf(origin).Elem(),
		reflect.TypeOf(local).Elem(),
		reflect.TypeOf(upstream).Elem(),
	).(reflect.Type)
	return mapStringAnyToPtrStruct(merged.(map[string]any), targetType)
}

// TODO: make these changes in twoway also
func threeWayMergeStruct(origin, local, upstream any) any {
	originMap := structToMap(origin)
	localMap := structToMap(local)
	upstreamMap := structToMap(upstream)

	merged := threeWayMergeMap(originMap, localMap, upstreamMap)
	if merged == nil {
		return reflect.Zero(reflect.TypeOf(local)).Interface()
	}

	targetType := threeWayMergeScalar(
		reflect.TypeOf(origin),
		reflect.TypeOf(local),
		reflect.TypeOf(upstream),
	).(reflect.Type)
	return mapStringAnyToStruct(merged.(map[string]any), targetType)
}

func threeWayMergeStructSlice(origin, local, upstream []any) any {
	if len(origin) == 0 {
		return twoWayMergeStructSlice(local, upstream)
	}

	if len(origin) == 0 && len(local) == 0 && len(upstream) == 0 {
		return threeWayMergeScalar(origin, local, upstream)
	}

	var targetType reflect.Type
	if len(origin) > 0 {
		targetType = reflect.TypeOf(origin[0])
	}

	if targetType.Kind() != reflect.Struct {
		return threeWayMergeSliceNonAssociative(origin, local, upstream)
	}

	if len(local) > 0 && reflect.TypeOf(local[0]) != targetType || len(upstream) > 0 && reflect.TypeOf(upstream[0]) != targetType {
		return threeWayMergeSliceNonAssociative(origin, local, upstream)
	}

	originSliceMap := sliceStructToSliceMap(origin)
	localSliceMap := sliceStructToSliceMap(local)
	upstreamSliceMap := sliceStructToSliceMap(upstream)

	mergedSliceMap := threeWayMergeSlice(originSliceMap, localSliceMap, upstreamSliceMap).([]any)
	return sliceMapToSliceStruct(mergedSliceMap, targetType)
}

func threeWayMergePtrStructSlice(origin, local, upstream []any) any {
	if len(origin) == 0 && len(local) == 0 && len(upstream) == 0 {
		return threeWayMergeScalar(origin, local, upstream)
	}

	if len(origin) == 0 {
		return twoWayMergePtrStructSlice(local, upstream)
	}

	var targetType reflect.Type
	if len(origin) > 0 {
		targetType = reflect.TypeOf(origin[0])
	}

	if targetType.Kind() != reflect.Ptr || targetType.Elem().Kind() != reflect.Struct {
		return threeWayMergeSliceNonAssociative(origin, local, upstream)
	}

	if len(local) > 0 && reflect.TypeOf(local[0]) != targetType || len(upstream) > 0 && reflect.TypeOf(upstream[0]) != targetType {
		return threeWayMergeSliceNonAssociative(origin, local, upstream)
	}

	originSliceMap := slicePtrStructToSliceMap(origin)
	localSliceMap := slicePtrStructToSliceMap(local)
	upstreamSliceMap := slicePtrStructToSliceMap(upstream)

	mergedSliceMap := threeWayMergeSlice(originSliceMap, localSliceMap, upstreamSliceMap).([]any)
	return sliceMapToSlicePtrStruct(mergedSliceMap, targetType)
}
