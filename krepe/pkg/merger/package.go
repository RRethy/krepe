package merger

import (
	"github.com/RRethy/krepe/krepe/pkg/pkg"
)

var _ Merger = Package{}

type Package struct{}

func (p Package) Merge(origin, local, upstream *pkg.Package) *pkg.Package {
	return threeWayMerge(origin, local, upstream).(*pkg.Package)
}
