package merger

type Merger[T any] interface {
	Merge(origin, local, upstream T) (T, error)
}
