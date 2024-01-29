package run

import (
	"github.com/RRethy/krepe/krepe/pkg/pkg"
)

type pipeline struct {
	pkg  *pkg.Package
	name string
}

func newPipeline(pkg *pkg.Package, name string) *pipeline {
	return nil
	// return &pipeline{pkg: pkg, name: name}
}

func (p *pipeline) run(dir string) error {
	// err := p.pkg.RunPipelineByName(p.name)
	// if err != nil {
	// 	return err
	// }
	//
	// err = p.pkg.Write(dir)
	// if err != nil {
	// 	return err
	// }
	//
	return nil
}
