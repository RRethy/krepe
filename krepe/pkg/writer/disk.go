package writer

import (
	"os"
	"path/filepath"

	"github.com/RRethy/krepe/krepe/pkg/pkg"
	"github.com/RRethy/krepe/krepe/pkg/yaml"
)

var _ Writer = Disk{}

type Disk struct {
	DirSuffix string
}

func (d Disk) Write(pkg *pkg.Package, dir string) error {
	if pkg == nil {
		return nil
	}

	pkgPath := filepath.Join(dir, pkg.Name, d.DirSuffix)
	err := os.MkdirAll(pkgPath, 0o755)
	if err != nil {
		return err
	}

	krepe := pkg.GetTypesKrepe()
	raw, err := yaml.Marshal(krepe)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(pkgPath, d.DirSuffix, "krepe.yaml"), raw, 0o644)
	if err != nil {
		return err
	}

	for _, fileImport := range pkg.FileImports {
		raw, err := yaml.Marshal(fileImport.Resource.Object)
		if err != nil {
			return err
		}

		err = os.WriteFile(filepath.Join(pkgPath, d.DirSuffix, fileImport.Name), raw, 0o644)
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
