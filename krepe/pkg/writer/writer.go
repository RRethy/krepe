package writer

import (
	"errors"

	"github.com/RRethy/krepe/krepe/pkg/pkg"
)

type Writer interface {
	Write(pkg *pkg.Package) error
}

type PackageWriter struct {
	dir string
}

func NewPackageWriter(dir string) (PackageWriter, error) {
	if dir == "" {
		return PackageWriter{}, errors.New("dir cannot be empty")
	}
	return PackageWriter{dir: dir}, nil
}

func (p PackageWriter) Write(pkg *pkg.Package) error {
	// TODO: write to disk
	return nil
}

type Noop struct{}

func (n Noop) Write(pkg *pkg.Package) error {
	return nil
}
