package threewaymerge

// Merge returns the result of performing a 3-way merge on the origin, local, and upstream resources.
func Merge(origin, local, upstream map[string]any) map[string]any {
	if res, ok := threeWayMergeMap(origin, local, upstream).(map[string]any); ok {
		return res
	}

	return nil
}
