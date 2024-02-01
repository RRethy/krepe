package merger

import (
	"github.com/RRethy/krepe/krepe/pkg/pkg"
)

var (
	_ Merger = Noop{}
	_ Merger = Package{}
)

type Merger interface {
	Merge(origin, local, upstream *pkg.Package) (*pkg.Package, error)
}

type Noop struct{}

func (n Noop) Merge(origin, local, upstream *pkg.Package) (*pkg.Package, error) {
	return local, nil
}

type Package struct{}

func NewMerger() Merger {
	return Package{}
}

func (p Package) Merge(origin, local, upstream *pkg.Package) (*pkg.Package, error) {
	return ThreeWayMerge(origin, local, upstream)
}
