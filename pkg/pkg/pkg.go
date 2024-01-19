package pkg

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Shopify/krepe/pkg/pkg/pipeline"
	"github.com/Shopify/krepe/pkg/pkg/resource"
)

type Pkg struct {
	name      string
	krepe     *Krepe
	resources []*resource.Resource
	packages  []*Pkg
}

func NewPkgFromPath(pkgPath string) (*Pkg, error) {
	if pkgPath == "/" {
		return nil, fmt.Errorf("pkg path cannot be `/`")
	}

	fileInfo, err := os.Stat(pkgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get info for pkg: %w", err)
	}
	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("pkg path is not a directory: %s", pkgPath)
	}

	name := filepath.Base(pkgPath)
	krepe, err := NewKrepeFromPath(filepath.Join(pkgPath, "krepe.yaml"))
	if err != nil {
		return nil, fmt.Errorf("failed to create krepe in pkg path `%s`: %w", pkgPath, err)
	}

	var resources []*resource.Resource
	for _, fileImport := range krepe.Imports.Files {
		resource, err := resource.NewResourceFromPath(filepath.Join(pkgPath, fileImport))
		if err != nil {
			return nil, fmt.Errorf("failed to create resource `%s` in pkg `%s`: %w", fileImport, pkgPath, err)
		}
		resources = append(resources, resource)
	}

	var packages []*Pkg
	for _, pkgImport := range krepe.Imports.Packages {
		pkg, err := NewPkgFromPath(filepath.Join(pkgPath, pkgImport.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed to create pkg `%s` in pkg `%s`: %w", pkgImport.Name(), pkgPath, err)
		}
		packages = append(packages, pkg)
	}

	return &Pkg{
		name:      name,
		krepe:     krepe,
		resources: resources,
		packages:  packages,
	}, nil
}

func (p *Pkg) RunPipelineByName(name string) error {
	if pipeline, ok := p.krepe.Pipelines[name]; ok {
		for _, resource := range p.resources {
			err := pipeline.Run(resource)
			if err != nil {
				return fmt.Errorf("failed to run pipeline `%s` on resource `%s`: %w", name, resource.Name, err)
			}
		}
		return nil
	} else {
		return fmt.Errorf("failed to get pipeline `%s`", name)
	}
}

func (p *Pkg) RunPipeline(pipeline pipeline.Pipeline) error {
	return nil
}
