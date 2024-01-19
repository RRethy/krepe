package merger

import (
	"github.com/RRethy/krepe/krepe/pkg/pkg"
)

var (
	_ Merger[*pkg.Pkg] = (*PkgMerger)(nil)
)

type PkgMerger struct {
}

func NewPkgMerger() *PkgMerger {
	return &PkgMerger{}
}

func (m *PkgMerger) TwoWayMerge(local, upstream *pkg.Pkg) *pkg.Pkg {
	if local == nil && upstream == nil {
		return nil
	}

	if local != nil && upstream == nil {
		return local
	}

	if local == nil && upstream != nil {
		return upstream
	}

	// name := local.Name
	// krepe := (&KrepeMerger{}).TwoWayMerge(local.Krepe, upstream.Krepe)

	panic("TODO")
}

func (m *PkgMerger) ThreeWayMerge(origin, local, upstream *pkg.Pkg) *pkg.Pkg {
	if origin == nil {
		return m.TwoWayMerge(local, upstream)
	}

	if local == nil && upstream == nil {
		return nil
	}

	if local == nil && upstream != nil {
		return nil
	}

	if local != nil && upstream == nil {
		return local
	}

	panic("TODO")
}
