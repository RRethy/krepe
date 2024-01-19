package merger

import (
	"github.com/RRethy/krepe/krepe/pkg/pkg"
)

var (
	_ Merger[*pkg.Krepe] = (*KrepeMerger)(nil)
)

type KrepeMerger struct {
}

func (m *KrepeMerger) TwoWayMerge(local, upstream *pkg.Krepe) *pkg.Krepe {
	panic("TODO")
}

func (m *KrepeMerger) ThreeWayMerge(origin, local, upstream *pkg.Krepe) *pkg.Krepe {
	panic("TODO")
}
