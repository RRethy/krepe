package threewaymerge

type mergeable interface {
	any | map[string]any | []any
}
