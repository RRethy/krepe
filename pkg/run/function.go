package run

import (
	"github.com/Shopify/krepe/pkg/pkg"
)

type function struct {
	pkg  *pkg.Pkg
	name string
}

func newFunction(pkg *pkg.Pkg, name string) *function {
	return &function{pkg: pkg, name: name}
}

func (f *function) run() error {
	return nil
}
