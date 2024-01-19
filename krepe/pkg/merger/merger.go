package merger

type Merger[T any] interface {
	TwoWayMerge(local, upstream T) T
	ThreeWayMerge(origin, local, upstream T) T
}
