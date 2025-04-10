package writer

import (
	"github.com/RRethy/krepe/krepe/pkg/pkg"
)

var _ Writer = Cache{}

type Cache struct{}

func (c Cache) Write(pkg *pkg.Package, dir string) error {
	return Disk{DirSuffix: ".krepe"}.Write(pkg, dir)
}
