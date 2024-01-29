package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/RRethy/krepe/krepe/pkg/git"
	"github.com/RRethy/krepe/krepe/pkg/types"
	"github.com/RRethy/krepe/krepe/pkg/yaml"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Package struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	PackageImports []PackageImport
	FileImports    []FileImport
	Pipelines      []Pipeline
}

func NewPackageFromPath(pkgPath string) (*Package, error) {
	return NewPackageFromPathWithName(pkgPath, filepath.Base(pkgPath))
}

func NewPackageFromPathWithName(packagePath, name string) (*Package, error) {
	if packagePath == "/" {
		return nil, fmt.Errorf("pkg path cannot be `/`")
	}

	if strings.Contains(name, "/") {
		return nil, fmt.Errorf("pkg name cannot have `/`")
	}

	// TODO: test this
	if name == "" {
		name = filepath.Base(packagePath)
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

		pkg, err := NewPackageFromPathWithName(filepath.Join(packagePath, name), name)
		if err != nil {
			return nil, fmt.Errorf("importing pkg `%s`: %w", name, err)
		}

		packageImports = append(packageImports, PackageImport{
			Ref:     pkgRef,
			Name:    name,
			Package: pkg,
		})
	}

	var fileImports []FileImport
	for _, fileImport := range krepe.Imports.Files {
		filename := filepath.Join(packagePath, fileImport)
		fileImport, err := NewFileImportFromPath(filename)
		if err != nil {
			return nil, fmt.Errorf("importing file `%s` in pkg `%s`: %w", fileImport.Filename, packagePath, err)
		}

		fileImports = append(fileImports, fileImport)
	}

	var pipelines []Pipeline
	for _, typesPipeline := range krepe.Pipelines {
		pipeline, err := NewPipeline(typesPipeline)
		if err != nil {
			return nil, fmt.Errorf("create pipeline `%s` in pkg `%s`: %w", typesPipeline.Name, packagePath, err)
		}

		pipelines = append(pipelines, pipeline)
	}

	krepe.ObjectMeta.Name = name

	return &Package{
		TypeMeta:       krepe.TypeMeta,
		ObjectMeta:     krepe.ObjectMeta,
		PackageImports: packageImports,
		FileImports:    fileImports,
		Pipelines:      pipelines,
	}, nil
}

func (p *Package) RunPipelineByName(name string) error {
	for _, pipeline := range p.Pipelines {
		if pipeline.Name == name {
			return p.RunPipeline(pipeline)
		}
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

func (p *Package) AddPackage(pkg *Package, ref *git.PkgRef, name string) error {
	if name == "" {
		name = ref.Name
	}

	for _, existingPkgImport := range p.PackageImports {
		if existingPkgImport.Name == name {
			return fmt.Errorf("pkg `%s` already exists", name)
		}
	}

	p.PackageImports = append(p.PackageImports, PackageImport{
		Ref:     ref,
		Name:    pkg.Name,
		Package: pkg,
	})
	return nil
}
