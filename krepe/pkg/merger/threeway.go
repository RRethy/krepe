package merger

import (
	"reflect"

	"github.com/RRethy/krepe/krepe/pkg/pkg"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	case *pkg.Package:
		upstreamPkg, ok := upstream.(*pkg.Package)
		if !ok {
			return local
		}

		originPkg, ok := origin.(*pkg.Package)
		if !ok {
			return twoWayMergePkg(localTyped, upstreamPkg)
		}

		return threeWayMergePkg(originPkg, localTyped, upstreamPkg)
	default:
		localType := reflect.TypeOf(local)
		upstreamType := reflect.TypeOf(upstream)
		originType := reflect.TypeOf(origin)

		if localType != upstreamType {
			return local
		}

		if localType != originType {
			return twoWayMerge(local, upstream)
		}

		if localType.Kind() == reflect.Ptr &&
			localType.Elem().Name() == upstreamType.Elem().Name() &&
			localType.Elem().Name() == originType.Elem().Name() {
			return threeWayMergePtrStruct(origin, local, upstream)
		} else if localType.Kind() == reflect.Struct &&
			localType.Name() == upstreamType.Name() &&
			localType.Name() == originType.Name() {
			return threeWayMergeStruct(origin, local, upstream)
		} else {
			return threeWayMergeScalar(origin, local, upstream)
		}
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
func threeWayMergeScalar[T any](origin, local, upstream T) T {
	if reflect.DeepEqual(origin, upstream) {
		return local
	}

	return upstream
}

func threeWayMergePkg(origin, local, upstream *pkg.Package) *pkg.Package {
	// PackageImports []PackageImport
	// FileImports    []FileImport
	// Pipelines      []Pipeline

	return &pkg.Package{
		TypeMeta: metav1.TypeMeta{
			Kind:       threeWayMergeScalar(origin.TypeMeta.Kind, local.TypeMeta.Kind, upstream.TypeMeta.Kind),
			APIVersion: threeWayMergeScalar(origin.TypeMeta.APIVersion, local.TypeMeta.APIVersion, upstream.TypeMeta.APIVersion),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:                       threeWayMergeScalar(origin.ObjectMeta.Name, local.ObjectMeta.Name, upstream.ObjectMeta.Name),
			GenerateName:               threeWayMergeScalar(origin.ObjectMeta.GenerateName, local.ObjectMeta.GenerateName, upstream.ObjectMeta.GenerateName),
			Namespace:                  threeWayMergeScalar(origin.ObjectMeta.Namespace, local.ObjectMeta.Namespace, upstream.ObjectMeta.Namespace),
			UID:                        threeWayMergeScalar(origin.ObjectMeta.UID, local.ObjectMeta.UID, upstream.ObjectMeta.UID),
			ResourceVersion:            threeWayMergeScalar(origin.ObjectMeta.ResourceVersion, local.ObjectMeta.ResourceVersion, upstream.ObjectMeta.ResourceVersion),
			Generation:                 threeWayMergeScalar(origin.ObjectMeta.Generation, local.ObjectMeta.Generation, upstream.ObjectMeta.Generation),
			CreationTimestamp:          threeWayMergeScalar(origin.ObjectMeta.CreationTimestamp, local.ObjectMeta.CreationTimestamp, upstream.ObjectMeta.CreationTimestamp),
			DeletionTimestamp:          threeWayMergeScalar(origin.ObjectMeta.DeletionTimestamp, local.ObjectMeta.DeletionTimestamp, upstream.ObjectMeta.DeletionTimestamp),
			DeletionGracePeriodSeconds: threeWayMergeScalar(origin.ObjectMeta.DeletionGracePeriodSeconds, local.ObjectMeta.DeletionGracePeriodSeconds, upstream.ObjectMeta.DeletionGracePeriodSeconds),
			Labels:                     threeWayMergeMapStringString(origin.ObjectMeta.Labels, local.ObjectMeta.Labels, upstream.ObjectMeta.Labels),
			Annotations:                threeWayMergeMapStringString(origin.ObjectMeta.Annotations, local.ObjectMeta.Annotations, upstream.ObjectMeta.Annotations),
			OwnerReferences:            threeWayMergeSliceStruct(origin.ObjectMeta.OwnerReferences, local.ObjectMeta.OwnerReferences, upstream.ObjectMeta.OwnerReferences),
			Finalizers:                 threeWayMergeSliceString(origin.ObjectMeta.Finalizers, local.ObjectMeta.Finalizers, upstream.ObjectMeta.Finalizers),
			ManagedFields:              threeWayMergeSliceStruct(origin.ObjectMeta.ManagedFields, local.ObjectMeta.ManagedFields, upstream.ObjectMeta.ManagedFields),
		},
		PackageImports: nil,
		FileImports:    nil,
		Pipelines:      nil,
	}
}

func threeWayMergeMapStringString(origin, local, upstream map[string]string) map[string]string {
	originMapStringAny, localMapStringAny, upstreamMapStringAny := make(map[string]any), make(map[string]any), make(map[string]any)
	for k, v := range origin {
		originMapStringAny[k] = any(v)
	}
	for k, v := range local {
		localMapStringAny[k] = any(v)
	}
	for k, v := range upstream {
		upstreamMapStringAny[k] = any(v)
	}
	labels := threeWayMergeMap(originMapStringAny, localMapStringAny, upstreamMapStringAny)
	labelsTyped := make(map[string]string)
	for k, v := range labels {
		labelsTyped[k] = v.(string)
	}
	return labelsTyped
}

func threeWayMergeSliceString(origin, local, upstream []string) []string {
	originSliceAny, localSliceAny, upstreamSliceAny := make([]any, len(origin)), make([]any, len(local)), make([]any, len(upstream))
	for i, v := range origin {
		originSliceAny[i] = any(v)
	}
	for i, v := range local {
		localSliceAny[i] = any(v)
	}
	for i, v := range upstream {
		upstreamSliceAny[i] = any(v)
	}
	finalizers := threeWayMergeSlice(originSliceAny, localSliceAny, upstreamSliceAny)
	finalizersTyped := make([]string, len(finalizers))
	for i, v := range finalizers {
		finalizersTyped[i] = v.(string)
	}
	return finalizersTyped
}

func threeWayMergeSliceStruct[T any](origin, local, upstream []T) []T {
	originSliceAny, localSliceAny, upstreamSliceAny := make([]any, len(origin)), make([]any, len(local)), make([]any, len(upstream))
	for i, v := range origin {
		originSliceAny[i] = any(structToMap(v))
	}

	for i, v := range local {
		localSliceAny[i] = any(structToMap(v))
	}

	for i, v := range upstream {
		upstreamSliceAny[i] = any(structToMap(v))
	}

	res := threeWayMergeSlice(originSliceAny, localSliceAny, upstreamSliceAny)
	resTyped := make([]T, len(res))
	for i, v := range res {
		m, ok := v.(map[string]any)
		if !ok {
			continue
		}
		resTyped[i] = mapStringAnyToStruct[T](m)
	}

	return resTyped
}

func threeWayMergePtrStruct(origin, local, upstream any) any {
	return local
}

func threeWayMergeStruct(origin, local, upstream any) any {
	return local
}
