package merger

import (
	"github.com/RRethy/krepe/krepe/pkg/pkg"
)

type Merger interface {
	Merge(origin, local, upstream *pkg.Package) *pkg.Package
}
