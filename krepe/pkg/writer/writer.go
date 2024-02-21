package writer

import (
	"github.com/RRethy/krepe/krepe/pkg/pkg"
)

type Writer interface {
	Write(pkg *pkg.Package, dir string) error
}
