package run

import (
	"github.com/Shopify/krepe/cli/pkg/pkg"
)

type pipeline struct {
	pkg  *pkg.Pkg
	name string
}

func newPipeline(pkg *pkg.Pkg, name string) *pipeline {
	return &pipeline{pkg: pkg, name: name}
}

func (p *pipeline) run() error {
	return nil
	// return p.pkg.RunPipeline(p.name)
}
