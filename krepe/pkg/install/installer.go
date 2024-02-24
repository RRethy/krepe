package install

import (
	"path/filepath"

	"github.com/RRethy/krepe/krepe/pkg/pkg"
	"github.com/RRethy/krepe/krepe/pkg/writer"
)

type Installer struct {
	Writer writer.Writer
}

func (installer *Installer) Install(pkgPath, newPkgPath, newPkgName string) error {
	rootPkg, err := pkg.NewPackageFromPath(pkgPath)
	if err != nil {
		return err
	}

	newPkgPath = filepath.Join(pkgPath, newPkgPath)
	newPkg, err := pkg.NewPackageFromPathWithName(newPkgPath, newPkgName)
	if err != nil {
		return err
	}

	relPath, err := filepath.Rel(pkgPath, newPkgPath)
	if err != nil {
		return err
	}

	err = rootPkg.AddPackage(newPkg, relPath)
	if err != nil {
		return err
	}

	err = installer.Writer.Write(rootPkg, pkgPath)
	if err != nil {
		return err
	}

	return nil
}
