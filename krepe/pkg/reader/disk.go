package reader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/RRethy/krepe/krepe/pkg/pkg"
	"github.com/RRethy/krepe/krepe/pkg/types"
	"github.com/RRethy/krepe/krepe/pkg/yaml"
)

var (
	_ Reader = Disk{}
)

type Disk struct {
	DirSuffix                  string
	AllowMissingPackageImports bool
}

func (d Disk) Read(pkgPath string) (*pkg.Package, error) {
	_, err := os.Stat(filepath.Join(pkgPath, d.DirSuffix))
	if err != nil {
		return nil, err
	}

	krepePath := filepath.Join(pkgPath, d.DirSuffix, "krepe.yaml")
	krepeYaml, err := os.ReadFile(krepePath)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", krepePath, err)
	}

	var krepe types.Krepe
	err = yaml.Unmarshal(krepeYaml, &krepe, yaml.DisallowUnknownFieldOption)
	if err != nil {
		return nil, fmt.Errorf("unmarshal %s: %w", krepePath, err)
	}

	var packageImports []pkg.PackageImport
	for _, packageImport := range krepe.Imports.Packages {
		subPkg, err := d.Read(filepath.Join(pkgPath, packageImport.Name))
		if err != nil {
			if !d.AllowMissingPackageImports {
				return nil, err
			}
			continue
		}

		subPkg.Name = packageImport.Name

		packageImports = append(packageImports, pkg.PackageImport{
			RelativePath: packageImport.RelativePath,
			Package:      subPkg,
		})
	}

	var fileImports []pkg.FileImport
	for _, fileImport := range krepe.Imports.Files {
		filename := filepath.Join(pkgPath, d.DirSuffix, fileImport)
		fileImport, err := pkg.NewFileImportFromPath(filename)
		if err != nil {
			return nil, err
		}

		fileImports = append(fileImports, fileImport)
	}

	var pipelines []pkg.Pipeline
	for _, typesPipeline := range krepe.Pipelines {
		pipeline, err := pkg.NewPipeline(typesPipeline)
		if err != nil {
			return nil, fmt.Errorf("create pipeline `%s` in pkg `%s`: %w", typesPipeline.Name, pkgPath, err)
		}

		pipelines = append(pipelines, pipeline)
	}

	return &pkg.Package{
		TypeMeta:       krepe.TypeMeta,
		ObjectMeta:     krepe.ObjectMeta,
		PackageImports: packageImports,
		FileImports:    fileImports,
		Pipelines:      pipelines,
	}, nil
}
