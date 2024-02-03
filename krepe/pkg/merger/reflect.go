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

func mapStringAnyToStruct[T any](m map[string]any) T {
	return *mapStringAnyToPtrStruct[T](m)
}

func mapStringAnyToPtrStruct[T any](m map[string]any) *T {
	var obj T
	typeOf := reflect.TypeOf(&obj).Elem()
	valueOf := reflect.ValueOf(&obj).Elem()
	for i := 0; i < typeOf.NumField(); i++ {
		field := valueOf.Field(i)
		name := typeOf.Field(i).Name
		if value, ok := m[name]; ok {
			field.Set(reflect.ValueOf(value))
		}
	}
	return &obj
}
