package run

import (
	"github.com/RRethy/krepe/krepe/pkg/pkg"
	"github.com/RRethy/krepe/krepe/pkg/writer"
)

type Pipeline struct {
	Pkg    *pkg.Package
	Name   string
	Writer writer.Writer
	Dir    string
}

func (p *Pipeline) Run() error {
	err := p.Pkg.RunPipelineByName(p.Name)
	if err != nil {
		return err
	}

	err = p.Writer.Write(p.Pkg, p.Dir)
	if err != nil {
		return err
	}

	return nil
}
