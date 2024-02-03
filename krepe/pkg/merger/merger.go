package merger

import (
	"github.com/RRethy/krepe/krepe/pkg/pkg"
)

var (
	_ Merger = Noop{}
	_ Merger = Package{}
)

type Merger interface {
	Merge(origin, local, upstream *pkg.Package) *pkg.Package
}

type Noop struct{}

func (n Noop) Merge(origin, local, upstream *pkg.Package) *pkg.Package {
	return local
}

type Package struct{}

func NewMerger() Merger {
	return Package{}
}

func (p Package) Merge(origin, local, upstream *pkg.Package) *pkg.Package {
	return threeWayMerge(origin, local, upstream).(*pkg.Package)
}
