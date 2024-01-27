package pkg

import (
	"github.com/RRethy/krepe/krepe/pkg/git"
)

// TODO: validate package.name == name
type PackageImport struct {
	Ref     *git.PkgRef
	Name    string
	Package *Package
}
