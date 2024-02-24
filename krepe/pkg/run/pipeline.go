package run

import (
	"github.com/RRethy/krepe/krepe/pkg/pkg"
	"github.com/RRethy/krepe/krepe/pkg/writer"
)

type Pipeline struct {
	Writer writer.Writer
}

func (p *Pipeline) Run(pkgPath, pipelineName string) error {
	pkg, err := pkg.NewPackageFromPath(pkgPath)
	if err != nil {
		return err
	}

	err = pkg.RunPipelineByName(pipelineName)
	if err != nil {
		return err
	}

	err = p.Writer.Write(pkg, pkgPath)
	if err != nil {
		return err
	}

	return nil
}
