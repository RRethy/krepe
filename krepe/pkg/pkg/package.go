package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	return NewPackageFromPathWithName(pkgPath, "")
}

func NewPackageFromPathWithName(pkgPath, name string) (*Package, error) {
	absPkgPath, err := filepath.Abs(pkgPath)
	if err != nil {
		return nil, err
	}

	if name == "" {
		name = filepath.Base(absPkgPath)
	}

	if name == "" || name == "." || name == "/" || strings.Contains(name, "/") {
		return nil, fmt.Errorf("pkg name cannot be `%s`", name)
	}

	fileInfo, err := os.Stat(absPkgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get info for pkg: %w", err)
	}
	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("pkg path is not a directory: %s", absPkgPath)
	}

	krepePath := filepath.Join(absPkgPath, "krepe.yaml")
	krepeYaml, err := os.ReadFile(krepePath)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", krepePath, err)
	}

	var krepe types.Krepe
	err = yaml.Unmarshal(krepeYaml, &krepe, yaml.DisallowUnknownFieldOption)
	if err != nil {
		return nil, fmt.Errorf("unmarshal %s: %w", krepePath, err)
	}

	var packageImports []PackageImport
	for _, packageImport := range krepe.Imports.Packages {
		pkg, err := NewPackageFromPath(filepath.Join(absPkgPath, packageImport.Name))
		if err != nil {
			return nil, fmt.Errorf("importing pkg `%s`: %w", name, err)
		}

		packageImports = append(packageImports, PackageImport{
			RelativePath: packageImport.RelativePath,
			Package:      pkg,
		})
	}

	var fileImports []FileImport
	for _, fileImport := range krepe.Imports.Files {
		filename := filepath.Join(pkgPath, fileImport)
		fileImport, err := NewFileImportFromPath(filename)
		if err != nil {
			return nil, fmt.Errorf("importing file `%s` in pkg `%s`: %w", fileImport.Name, pkgPath, err)
		}

		fileImports = append(fileImports, fileImport)
	}

	var pipelines []Pipeline
	for _, typesPipeline := range krepe.Pipelines {
		pipeline, err := NewPipeline(typesPipeline)
		if err != nil {
			return nil, fmt.Errorf("create pipeline `%s` in pkg `%s`: %w", typesPipeline.Name, pkgPath, err)
		}

		pipelines = append(pipelines, pipeline)
	}

	krepe.Name = name

	return &Package{
		TypeMeta:       krepe.TypeMeta,
		ObjectMeta:     krepe.ObjectMeta,
		PackageImports: packageImports,
		FileImports:    fileImports,
		Pipelines:      pipelines,
	}, nil
}

func (p *Package) RunPipelineByName(name string) (found bool, err error) {
	found = false
	for _, packageImport := range p.PackageImports {
		foundInImport, err := packageImport.Package.RunPipelineByName(name)
		if err != nil {
			return false, err
		}
		found = found || foundInImport
	}

	for _, pipeline := range p.Pipelines {
		if pipeline.Name == name {
			return true, p.RunPipeline(pipeline)
		}
	}

	return found, nil
}

func (p *Package) RunPipeline(pipeline Pipeline) error {
	for _, fileImport := range p.FileImports {
		err := pipeline.Run(fileImport.Resource)
		if err != nil {
			return fmt.Errorf("running pipeline on resource `%s`: %w", fileImport.Name, err)
		}
	}

	for _, pkgImport := range p.PackageImports {
		err := pkgImport.Package.RunPipeline(pipeline)
		if err != nil {
			return fmt.Errorf("running pipeline on pkg `%s`: %w", pkgImport.Package.Name, err)
		}
	}

	return nil
}

func (p *Package) AddPackage(pkg *Package, relPath string) error {
	for _, existingPkgImport := range p.PackageImports {
		if existingPkgImport.Package.Name == pkg.Name {
			return fmt.Errorf("package with name `%s` already exists", pkg.Name)
		}
	}

	p.PackageImports = append(p.PackageImports, PackageImport{
		RelativePath: relPath,
		Package:      pkg,
	})
	return nil
}

func (p *Package) GetPackageImportByName(name string) (pkgImport *PackageImport, ok bool) {
	for _, pkgImport := range p.PackageImports {
		if pkgImport.Package.Name == name {
			return &pkgImport, true
		}
	}

	return nil, false
}

func (p *Package) UpdatePackage(pkg *Package, relPath string) error {
	for i, existingPkgImport := range p.PackageImports {
		if existingPkgImport.Package.Name == pkg.Name {
			p.PackageImports[i] = PackageImport{
				RelativePath: relPath,
				Package:      pkg,
			}
			return nil
		}
	}

	return fmt.Errorf("no imported package with name %s", pkg.Name)
}

func (p *Package) GetTypesKrepe() *types.Krepe {
	var fileImports []string
	for _, fileImport := range p.FileImports {
		fileImports = append(fileImports, fileImport.Name)
	}

	var packageImports []types.PackageImport
	for _, packageImport := range p.PackageImports {
		packageImports = append(packageImports, types.PackageImport{
			RelativePath: packageImport.RelativePath,
			Name:         packageImport.Package.Name,
		})
	}

	var pipelines []types.Pipeline
	for _, pipeline := range p.Pipelines {
		var steps []types.Step
		for _, step := range pipeline.Steps {
			steps = append(steps, types.Step{
				Function: step.Name,
				Target:   step.Target.ToTypesTarget(),
				Config:   step.Config,
			})
		}

		pipelines = append(pipelines, types.Pipeline{
			Name:  pipeline.Name,
			Steps: steps,
		})
	}

	return &types.Krepe{
		TypeMeta:   p.TypeMeta,
		ObjectMeta: p.ObjectMeta,
		Imports: types.Imports{
			Packages: packageImports,
			Files:    fileImports,
		},
		Pipelines: pipelines,
	}
}
