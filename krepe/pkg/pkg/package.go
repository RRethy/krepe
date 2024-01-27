package pkg

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/RRethy/krepe/krepe/pkg/git"
	"github.com/RRethy/krepe/krepe/pkg/types"
	"github.com/RRethy/krepe/krepe/pkg/yaml"
	"github.com/wk8/go-ordered-map/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Package struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	PackageImports []PackageImport
	FileImports    []FileImport
	Pipelines      *orderedmap.OrderedMap[string, Pipeline]
}

func NewPackageFromPath(pkgPath string) (*Package, error) {
	return NewPackageFromPathWithName(pkgPath, filepath.Base(pkgPath))
}

func NewPackageFromPathWithName(packagePath, name string) (*Package, error) {
	if packagePath == "/" {
		return nil, fmt.Errorf("pkg path cannot be `/`")
	}

	fileInfo, err := os.Stat(packagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get info for pkg: %w", err)
	}
	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("pkg path is not a directory: %s", packagePath)
	}

	krepePath := filepath.Join(packagePath, "krepe.yaml")
	krepeYaml, err := os.ReadFile(krepePath)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", krepePath, err)
	}

	var krepe types.Krepe
	err = yaml.Unmarshal(krepeYaml, &krepe)
	if err != nil {
		return nil, fmt.Errorf("unmarshal %s: %w", krepePath, err)
	}

	var packageImports []PackageImport
	for _, packageImport := range krepe.Imports.Packages {
		pkgRef, err := git.NewPkgRefFromString(packageImport.Ref)
		if err != nil {
			return nil, fmt.Errorf("parse package import ref `%s`: %w", packageImport.Ref, err)
		}

		name := pkgRef.Name
		if packageImport.Name != "" {
			name = packageImport.Name
		}

		packageImports = append(packageImports, PackageImport{
			Ref:     pkgRef,
			Name:    name,
			Package: nil,
		})
	}

	var fileImports []FileImport
	for _, fileImport := range krepe.Imports.Files {
		filename := filepath.Join(packagePath, fileImport)
		resource, err := NewResourceFromPath(filename)
		if err != nil {
			return nil, fmt.Errorf("create resource `%s` in pkg `%s`: %w", fileImport, packagePath, err)
		}

		fileImports = append(fileImports, FileImport{
			Filename: filename,
			Resource: resource,
		})
	}

	pipelines := orderedmap.New[string, Pipeline]()
	for pair := krepe.Pipelines.Oldest(); pair != nil; pair = pair.Next() {
		pipeline, err := NewPipeline(pair.Value)
		if err != nil {
			return nil, fmt.Errorf("create pipeline `%s` in pkg `%s`: %w", pair.Key, packagePath, err)
		}

		pipelines.Set(pair.Key, pipeline)
	}

	return &Package{
		TypeMeta:       krepe.TypeMeta,
		ObjectMeta:     krepe.ObjectMeta,
		PackageImports: packageImports,
		FileImports:    fileImports,
		Pipelines:      pipelines,
	}, nil
}

func (p *Package) RunPipelineByName(name string) error {
	if pipeline, ok := p.Pipelines.Get(name); ok {
		return p.RunPipeline(pipeline)

	}
	return fmt.Errorf("failed to get pipeline `%s`", name)
}

func (p *Package) RunPipeline(pipeline Pipeline) error {
	for _, fileImport := range p.FileImports {
		err := pipeline.Run(fileImport.Resource)
		if err != nil {
			return fmt.Errorf("running pipeline on resource `%s`: %w", fileImport.Filename, err)
		}
	}

	for _, pkgImport := range p.PackageImports {
		err := pkgImport.Package.RunPipeline(pipeline)
		if err != nil {
			return fmt.Errorf("running pipeline on pkg `%s`: %w", pkgImport.Name, err)
		}
	}

	return nil
}

func (p *Package) Write(dir string) error {
	return nil
}

// func (p *Pkg) Write(dir string) error {
// 	pkgPath := filepath.Join(dir, p.Name)
//
// 	err := os.MkdirAll(pkgPath, 0755)
// 	if err != nil {
// 		return fmt.Errorf("failed to create pkg directory `%s`: %w", pkgPath, err)
// 	}
//
// 	err = p.Krepe.Write(pkgPath)
// 	if err != nil {
// 		return fmt.Errorf("failed to write krepe.yaml in pkg directory `%s`: %w", pkgPath, err)
// 	}
//
// 	for _, resource := range p.resources {
// 		err = resource.Write(pkgPath)
// 		if err != nil {
// 			return fmt.Errorf("failed to write resource `%s` in pkg directory `%s`: %w", resource.Fname(), pkgPath, err)
// 		}
// 	}
//
// 	for _, pkg := range p.packages {
// 		err = pkg.Write(pkgPath)
// 		if err != nil {
// 			return fmt.Errorf("failed to write pkg `%s` in pkg directory `%s`: %w", pkg.Name, pkgPath, err)
// 		}
// 	}
//
// 	return nil
// }
