package run

import (
	"github.com/RRethy/krepe/krepe/pkg/pkg"
)

type pipeline struct {
	pkg  *pkg.Pkg
	name string
}

func newPipeline(pkg *pkg.Pkg, name string) *pipeline {
	return &pipeline{pkg: pkg, name: name}
}

func (p *pipeline) run(dir string) error {
	err := p.pkg.RunPipelineByName(p.name)
	if err != nil {
		return err
	}

	err = p.pkg.Write(dir)
	if err != nil {
		return err
	}

	return nil
}
