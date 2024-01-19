package krmmerge

var (
	// TODO: this is a bit of a hack, but it works for now, we should try
	// to use openapi
	associativeKeys = []string{
		"mountPath",
		"devicePath",
		"ip",
		"type",
		"topologyKey",
		"name",
		"containerPort",
	}
)

// isAssociativeSlice returns true if the slice contains only associations,
// otherwise false.
func isAssociativeSlice(slice []any) bool {
	if len(slice) == 0 {
		// it doesn't actually matter what we return here
		return true
	}

	for _, v := range slice {
		switch v.(type) {
		case map[string]any:
		default:
			return false
		}
	}

	return true
}

func getAssociativeKeys(slice []any) []string {
	if len(slice) == 0 {
		return nil
	}

	counts := make(map[string]map[string]struct{}, len(associativeKeys))
	for _, key := range associativeKeys {
		counts[key] = make(map[string]struct{})

		for _, elem := range slice {
			elemMap, ok := elem.(map[string]any)
			if !ok {
				return nil
			}

			keyVal, ok := elemMap[key].(string)
			if !ok {
				break
			}

			counts[key][keyVal] = struct{}{}
		}
	}

	var res []string
	for _, key := range associativeKeys {
		if len(counts[key]) == len(slice) {
			res = append(res, key)
		}
	}

	return res
}

func hasAssociativeKey(slice []any, key string) bool {
	values := make(map[string]struct{}, len(slice))

	for _, elem := range slice {
		elemMap, ok := elem.(map[string]any)
		if !ok {
			return false
		}

		val, ok := elemMap[key]
		if !ok {
			return false
		}

		valStr, ok := val.(string)
		if !ok {
			return false
		}

		values[valStr] = struct{}{}
	}

	return len(values) == len(slice)
}

func getCommonAssociativeKey(keyss ...[]string) string {
	commonKeys := getCommonAssociativeKeys(keyss)
	if len(commonKeys) == 0 {
		return ""
	}

	return commonKeys[0]
}

func getCommonAssociativeKeys(keyss [][]string) []string {
	if len(keyss) == 0 {
		return nil
	}

	if len(keyss) == 1 {
		return keyss[0]
	}

	if len(keyss) > 2 {
		return getCommonAssociativeKeys([][]string{
			keyss[0],
			getCommonAssociativeKeys(keyss[1:]),
		})
	}

	keys1 := keyss[0]
	keys2 := keyss[1]
	keys2Set := make(map[string]struct{}, len(keys2))
	for _, key := range keys2 {
		keys2Set[key] = struct{}{}
	}

	var res []string
	for _, key := range keys1 {
		if _, ok := keys2Set[key]; ok {
			res = append(res, key)
		}
	}

	return res
}
