package merger

import (
	"github.com/RRethy/krepe/krepe/pkg/pkg"
)

var (
	_ Merger = Noop{}
	_ Merger = PackageMerger{}
)

type Merger interface {
	Merge(origin, local, upstream *pkg.Package) (*pkg.Package, error)
}

type Noop struct{}

func (n Noop) Merge(origin, local, upstream *pkg.Package) (*pkg.Package, error) {
	return local, nil
}

type PackageMerger struct{}

func NewMerger() Merger {
	return PackageMerger{}
}

func (p PackageMerger) Merge(origin, local, upstream *pkg.Package) (*pkg.Package, error) {
	return ThreeWayMerge(origin, local, upstream)
}
