package writer

import (
	"errors"

	"github.com/RRethy/krepe/krepe/pkg/pkg"
)

type Writer interface {
	Write(pkg *pkg.Package) error
}

type Package struct {
	dir string
}

func NewPackageWriter(dir string) (Package, error) {
	if dir == "" {
		return Package{}, errors.New("dir cannot be empty")
	}
	return Package{dir: dir}, nil
}

func (p Package) Write(pkg *pkg.Package) error {
	// TODO: write to disk
	return nil
}

type Noop struct{}

func (n Noop) Write(pkg *pkg.Package) error {
	return nil
}
