package writer

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/RRethy/krepe/krepe/pkg/pkg"
	"github.com/RRethy/krepe/krepe/pkg/yaml"
)

type Disk struct {
}

func NewDiskWriter() (Disk, error) {
	return Disk{}, nil
}

func (d Disk) Write(pkg *pkg.Package, dir string) error {
	if dir == "" {
		return errors.New("dir cannot be empty")
	}

	if pkg == nil {
		return nil
	}

	pkgPath := filepath.Join(dir, pkg.Name)
	err := os.MkdirAll(pkgPath, 0755)
	if err != nil {
		return err
	}

	krepe := pkg.GetKrepe()
	raw, err := yaml.Marshal(krepe)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(pkgPath, "krepe.yaml"), raw, 0644)
	if err != nil {
		return err
	}

	for _, fileImport := range pkg.FileImports {
		raw, err := yaml.Marshal(fileImport.Resource)
		if err != nil {
			return err
		}

		err = os.WriteFile(filepath.Join(pkgPath, fileImport.Name), raw, 0644)
		if err != nil {
			return err
		}
	}

	for _, packageImport := range pkg.PackageImports {
		err = d.Write(packageImport.Package, filepath.Join(pkgPath, packageImport.Package.Name))
		if err != nil {
			return err
		}
	}

	return nil
}
