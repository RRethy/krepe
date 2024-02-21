package writer

import (
	"github.com/RRethy/krepe/krepe/pkg/pkg"
)

type Noop struct{}

func (n Noop) Write(pkg *pkg.Package, dir string) error {
	return nil
}
