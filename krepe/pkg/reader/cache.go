package reader

import (
	"github.com/RRethy/krepe/krepe/pkg/pkg"
)

var (
	_ Reader = Cache{}
)

type Cache struct{}

func (c Cache) Read(pkgPath string) (*pkg.Package, error) {
	return Disk{DirSuffix: ".krepe", AllowMissingPackageImports: true}.Read(pkgPath)
}
