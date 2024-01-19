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

func getAssociativeKey(slice []any) string {
	if len(slice) == 0 {
		return ""
	}

	counts := make(map[string]map[string]struct{}, len(associativeKeys))
	for _, key := range associativeKeys {
		counts[key] = make(map[string]struct{})

		for _, elem := range slice {
			elemMap, ok := elem.(map[string]any)
			if !ok {
				return ""
			}

			keyVal, ok := elemMap[key].(string)
			if !ok {
				break
			}

			counts[key][keyVal] = struct{}{}
		}
	}

	for _, key := range associativeKeys {
		if len(counts[key]) == len(slice) {
			return key
		}
	}

	return ""
}
