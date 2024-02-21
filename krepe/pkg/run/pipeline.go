package run

import (
	"github.com/RRethy/krepe/krepe/pkg/pkg"
	"github.com/RRethy/krepe/krepe/pkg/writer"
)

type pipeline struct {
	pkg    *pkg.Package
	name   string
	writer writer.Writer
	dir    string
}

type option func(*pipeline)

func WithWriter(w writer.Writer) option {
	return func(p *pipeline) {
		p.writer = w
	}
}

func WithDir(dir string) option {
	return func(p *pipeline) {
		p.dir = dir
	}
}

func newPipeline(pkg *pkg.Package, name string, options ...option) *pipeline {
	p := &pipeline{pkg: pkg, name: name, writer: writer.Noop{}}
	for _, o := range options {
		o(p)
	}

	return p
}

func (p *pipeline) run() error {
	err := p.pkg.RunPipelineByName(p.name)
	if err != nil {
		return err
	}

	err = p.writer.Write(p.pkg, p.dir)
	if err != nil {
		return err
	}

	return nil
}
