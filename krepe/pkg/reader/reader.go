package reader

import (
	"github.com/RRethy/krepe/krepe/pkg/pkg"
)

type Reader interface {
	Read(pkgPath string) (*pkg.Package, error)
}
