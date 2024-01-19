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
	counts := make(map[string]int, len(associativeKeys))
	for _, v := range slice {
		v, ok := v.(map[string]any)
		if !ok {
			return ""
		}

		for _, k := range associativeKeys {
			if _, ok := v[k]; ok {
				counts[k]++
			}
		}
	}

	for _, k := range associativeKeys {
		if counts[k] == len(slice) {
			return k
		}
	}

	return ""
}
