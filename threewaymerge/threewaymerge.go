// TODO: add docs

package threewaymerge

// Merge returns the result of performing a 3-way merge on the origin, local, and upstream maps.
// This method does not mutate origin, local, and upstream, but the result might share data structures.
func Merge(origin, local, upstream map[string]any) map[string]any {
	return threeWayMergeMap(origin, local, upstream)
}
