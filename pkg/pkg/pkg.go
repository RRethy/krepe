package pkg

import (
	"fmt"
	"os"
	"path/filepath"
)

type Pkg struct {
	name      string
	krepe     *Krepe
	resources []*resource
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

	var resources []*resource
	for _, fileImport := range krepe.Imports.Files {
		resource, err := newResourceFromPath(filepath.Join(pkgPath, fileImport))
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

func (p *Pkg) RunPipeline(name string) error {
	return nil
}
