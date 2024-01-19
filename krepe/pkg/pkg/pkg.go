package pkg

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Shopify/krepe/krepe/pkg/pkg/imports"
	"github.com/Shopify/krepe/krepe/pkg/pkg/resource"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type Pkg struct {
	name      string
	Krepe     *Krepe
	resources []*resource.Resource
	packages  []*Pkg
}

func NewPkgFromPathWithName(pkgPath, name string) (*Pkg, error) {
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

	krepe, err := NewKrepeFromPath(filepath.Join(pkgPath, "krepe.yaml"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse krepe.yaml in pkg path `%s`: %w", pkgPath, err)
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
		Krepe:     krepe,
		resources: resources,
		packages:  packages,
	}, nil
}

func NewPkgFromPath(pkgPath string) (*Pkg, error) {
	return NewPkgFromPathWithName(pkgPath, filepath.Base(pkgPath))
}

func (p *Pkg) RunPipelineByName(name string) error {
	if pipeline, ok := p.Krepe.Pipelines.Get(name); ok {
		for _, resource := range p.resources {
			err := pipeline.Run(resource)
			if err != nil {
				return fmt.Errorf("failed to run pipeline `%s` on resource `%s`: %w", name, resource.Fname(), err)
			}
		}
		return nil
	} else {
		return fmt.Errorf("failed to get pipeline `%s`", name)
	}
}

func (p *Pkg) InstallPackage(url string, name string) error {
	pkgImport, err := imports.NewPkg(url, name)
	if err != nil {
		return fmt.Errorf("creating pkg import `%s` from url `%s`: %w", name, url, err)
	}

	if p.PkgExists(pkgImport.Name()) {
		return fmt.Errorf("pkg `%s` already exists, consider giving it another name", pkgImport.Name())
	}

	// TODO: This is all temporary for a MVP
	tmpDir := "/tmp"
	pkgPath := filepath.Join(tmpDir, pkgImport.Name())

	if _, err := os.Stat(pkgPath); err == nil {
		if err := os.RemoveAll(pkgPath); err != nil {
			return fmt.Errorf("removing existing pkg `%s`: %w", pkgPath, err)
		}
	}

	repo, err := git.PlainClone(pkgPath, false, &git.CloneOptions{
		URL:      pkgImport.URL(),
		Progress: os.Stdout,
	})
	if err != nil {
		return fmt.Errorf("cloning pkg from url `%s` to `%s`: %w", pkgImport.URL(), pkgPath, err)
	}

	wt, err := repo.Worktree()
	if err != nil {
		return err
	}

	err = wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName("refs/tags/" + pkgImport.Version()),
	})
	if err != nil {
		return fmt.Errorf("checking out tag `%s` in pkg `%s`: %w", pkgImport.Version(), pkgPath, err)
	}

	pkg, err := NewPkgFromPathWithName(pkgPath, pkgImport.Name())
	if err != nil {
		return fmt.Errorf("creating pkg from path `%s`: %w", pkgPath, err)
	}

	p.packages = append(p.packages, pkg)
	p.Krepe.Imports.AddPackage(pkgImport)
	return nil
}

func (p *Pkg) Write(dir string) error {
	pkgPath := filepath.Join(dir, p.name)

	err := os.MkdirAll(pkgPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create pkg directory `%s`: %w", pkgPath, err)
	}

	err = p.Krepe.Write(pkgPath)
	if err != nil {
		return fmt.Errorf("failed to write krepe.yaml in pkg directory `%s`: %w", pkgPath, err)
	}

	for _, resource := range p.resources {
		err = resource.Write(pkgPath)
		if err != nil {
			return fmt.Errorf("failed to write resource `%s` in pkg directory `%s`: %w", resource.Fname(), pkgPath, err)
		}
	}

	for _, pkg := range p.packages {
		err = pkg.Write(pkgPath)
		if err != nil {
			return fmt.Errorf("failed to write pkg `%s` in pkg directory `%s`: %w", pkg.name, pkgPath, err)
		}
	}

	return nil
}

func (p *Pkg) PkgExists(other string) bool {
	for _, pkg := range p.packages {
		if pkg.name == other {
			return true
		}
	}

	return false
}
