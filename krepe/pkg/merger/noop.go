package merger

import (
	"github.com/RRethy/krepe/krepe/pkg/pkg"
)

var (
	_ Merger = Noop{}
)

type Noop struct{}

func (n Noop) Merge(origin, local, upstream *pkg.Package) *pkg.Package {
	return local
}
