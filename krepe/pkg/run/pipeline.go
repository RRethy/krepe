package run

import (
	"fmt"
	"path/filepath"

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

	found, err := pkg.RunPipelineByName(pipelineName)
	if !found {
		return fmt.Errorf("pipeline `%s` not found", pipelineName)
	}
	if err != nil {
		return err
	}

	err = p.Writer.Write(pkg, filepath.Dir(pkgPath))
	if err != nil {
		return err
	}

	return nil
}
