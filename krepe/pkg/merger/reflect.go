package merger

import (
	"reflect"
)

// TODO: we should panic occasionally
// TODO: handle ptrs being nil

func structToMap(obj any) map[string]any {
	if obj == nil || reflect.TypeOf(obj).Kind() != reflect.Struct {
		return make(map[string]any)
	}

	m := make(map[string]any)
	typeOf := reflect.TypeOf(obj)
	valueOf := reflect.ValueOf(obj)
	for i := 0; i < valueOf.NumField(); i++ {
		field := valueOf.Field(i)
		if field.CanInterface() {
			m[typeOf.Field(i).Name] = field.Interface()
		}
	}

	return m
}

func ptrStructToMap(obj any) map[string]any {
	if obj == nil ||
		reflect.TypeOf(obj).Kind() != reflect.Ptr ||
		reflect.TypeOf(obj).Elem().Kind() != reflect.Struct {
		return make(map[string]any)
	}

	m := make(map[string]any)
	typeOf := reflect.TypeOf(obj).Elem()
	valueOf := reflect.ValueOf(obj)
	if valueOf.IsNil() {
		return m
	}
	valueOf = valueOf.Elem()
	for i := 0; i < valueOf.NumField(); i++ {
		field := valueOf.Field(i)
		if field.CanInterface() {
			m[typeOf.Field(i).Name] = field.Interface()
		}
	}

	return m
}

func mapStringAnyToStruct(m map[string]any, structType reflect.Type) any {
	obj := reflect.New(structType).Elem()
	for i := 0; i < structType.NumField(); i++ {
		field := obj.Field(i)
		name := structType.Field(i).Name
		if value, ok := m[name]; ok {
			if field.CanSet() && value != nil {
				field.Set(reflect.ValueOf(value))
			}
		}
	}
	return obj.Interface()
}

func mapStringAnyToPtrStruct(m map[string]any, structType reflect.Type) any {
	if len(m) == 0 {
		return nil
	}

	if structType.Kind() == reflect.Pointer {
		structType = structType.Elem()
	}

	obj := reflect.New(structType)
	valueOf := obj.Elem()
	for i := 0; i < structType.NumField(); i++ {
		field := valueOf.Field(i)
		name := structType.Field(i).Name
		if value, ok := m[name]; ok {
			if value != nil {
				field.Set(reflect.ValueOf(value))
			}
		}
	}
	return obj.Interface()
}

func isUniformStructSlice(slice []any) bool {
	if len(slice) == 0 {
		return false
	}
	targetType := reflect.TypeOf(slice[0])
	if targetType.Kind() != reflect.Struct {
		return false
	}
	for _, item := range slice {
		if reflect.TypeOf(item) != targetType || reflect.TypeOf(item).Name() != targetType.Name() {
			return false
		}
	}
	return true
}

func isUniformPtrStructSlice(slice []any) bool {
	if len(slice) == 0 {
		return false
	}
	ptrType := reflect.TypeOf(slice[0])
	if ptrType.Kind() != reflect.Ptr {
		return false
	}
	targetType := ptrType.Elem()
	if targetType.Kind() != reflect.Struct {
		return false
	}
	for _, item := range slice {
		itemType := reflect.TypeOf(item)
		if itemType.Kind() != reflect.Ptr {
			return false
		}
		itemType = itemType.Elem()
		if itemType != targetType ||
			itemType.Name() != targetType.Name() {
			return false
		}
	}
	return true
}

func sliceStructToSliceMap(slice []any) []any {
	m := make([]any, len(slice))
	for i, item := range slice {
		m[i] = structToMap(item)
	}
	return m
}

func slicePtrStructToSliceMap(slice []any) []any {
	m := make([]any, len(slice))
	for i, item := range slice {
		m[i] = ptrStructToMap(item)
	}
	return m
}

func isUniformStructSlices(slices ...[]any) bool {
	if len(slices) == 0 || len(slices[0]) == 0 {
		return false
	}

	targetType := reflect.TypeOf(slices[0][0])

	for _, slice := range slices {
		if len(slice) == 0 || !isUniformStructSlice(slice) {
			return false
		}
		itemType := reflect.TypeOf(slice[0])
		if itemType != targetType || itemType.Name() != targetType.Name() {
			return false
		}
	}

	return true
}

func isUniformPtrStructSlices(slices ...[]any) bool {
	if len(slices) == 0 || len(slices[0]) == 0 {
		return false
	}

	targetType := reflect.TypeOf(slices[0][0])
	if targetType.Kind() != reflect.Ptr {
		return false
	}
	targetType = targetType.Elem()
	if targetType.Kind() != reflect.Struct {
		return false
	}

	for _, slice := range slices {
		if len(slice) == 0 || !isUniformPtrStructSlice(slice) {
			return false
		}
		itemType := reflect.TypeOf(slice[0]).Elem()
		if itemType != targetType || itemType.Name() != targetType.Name() {
			return false
		}
	}

	return true
}

func sliceMapToSliceStruct(slice []any, structType reflect.Type) []any {
	m := make([]any, len(slice))
	for i, item := range slice {
		m[i] = mapStringAnyToStruct(item.(map[string]any), structType)
	}
	return m
}

func sliceMapToSlicePtrStruct(slice []any, structType reflect.Type) []any {
	m := make([]any, len(slice))
	for i, item := range slice {
		m[i] = mapStringAnyToPtrStruct(item.(map[string]any), structType)
	}
	return m
}
