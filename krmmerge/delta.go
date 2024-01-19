package krmmerge

import (
	"reflect"
)

func delta(source, remove any) any {
	if reflect.TypeOf(source) != reflect.TypeOf(remove) {
		return source
	}

	switch source.(type) {
	case map[string]any:
		return deltaMap(source.(map[string]any), remove.(map[string]any))
	case []any:
		return deltaSlice(source.([]any), remove.([]any))
	default:
		if reflect.DeepEqual(source, remove) {
			return nil
		}
		return source
	}
}

func deltaMap(source, remove map[string]any) map[string]any {
	res := make(map[string]any)

	for k, v := range source {
		if removeVal, ok := remove[k]; ok {
			newVal := delta(v, removeVal)
			if newVal != nil {
				res[k] = newVal
			}
		} else {
			res[k] = v
		}
	}

	if len(res) == 0 {
		return nil
	}
	return res
}

func deltaSlice(source, remove []any) []any {
	if len(source) != len(remove) {
		return source
	}

	for i, v := range source {
		if !reflect.DeepEqual(v, remove[i]) {
			return source
		}
	}

	return nil
}
