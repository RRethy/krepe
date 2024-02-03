package merger

import (
	"reflect"
)

func structToMap(obj any) map[string]any {
	if reflect.TypeOf(obj).Kind() != reflect.Struct {
		return make(map[string]any)
	}

	m := make(map[string]any)
	typeOf := reflect.TypeOf(obj)
	valueOf := reflect.ValueOf(obj)
	for i := 0; i < valueOf.NumField(); i++ {
		field := valueOf.Field(i)
		m[typeOf.Field(i).Name] = field.Interface()
	}

	return m
}

func ptrStructToMap(obj any) map[string]any {
	if reflect.TypeOf(obj).Kind() != reflect.Ptr ||
		reflect.TypeOf(obj).Elem().Kind() != reflect.Struct {
		return make(map[string]any)
	}

	m := make(map[string]any)
	typeOf := reflect.TypeOf(obj).Elem()
	valueOf := reflect.ValueOf(obj).Elem()
	for i := 0; i < valueOf.NumField(); i++ {
		field := valueOf.Field(i)
		m[typeOf.Field(i).Name] = field.Interface()
	}

	return m
}

func mapStringAnyToStruct(m map[string]any, structType reflect.Type) any {
	obj := reflect.New(structType).Elem()
	for i := 0; i < structType.NumField(); i++ {
		field := obj.Field(i)
		name := structType.Field(i).Name
		if value, ok := m[name]; ok {
			field.Set(reflect.ValueOf(value))
		}
	}
	return obj.Interface()
}

func mapStringAnyToPtrStruct(m map[string]any, structType reflect.Type) any {
	if structType.Kind() == reflect.Pointer {
		structType = structType.Elem()
	}

	obj := reflect.New(structType)
	valueOf := obj.Elem()
	for i := 0; i < structType.NumField(); i++ {
		field := valueOf.Field(i)
		name := structType.Field(i).Name
		if value, ok := m[name]; ok {
			field.Set(reflect.ValueOf(value))
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

func sliceStructToSliceMap(slice []any) []map[string]any {
	m := make([]map[string]any, len(slice))
	for i, item := range slice {
		m[i] = structToMap(item)
	}
	return m
}

func slicePtrStructToSliceMap(slice []any) []map[string]any {
	m := make([]map[string]any, len(slice))
	for i, item := range slice {
		m[i] = ptrStructToMap(item)
	}
	return m
}
