package install

import (
	"fmt"
	"path/filepath"

	"github.com/RRethy/krepe/krepe/pkg/pkg"
)

func Install(pkgPath, url, name string) error {
	absPath, err := filepath.Abs(pkgPath)
	if err != nil {
		return fmt.Errorf("getting absolute path: %w", err)
	}

	p, err := pkg.NewPackageFromPath(absPath)
	if err != nil {
		return err
	}

	installer, err := NewInstaller()
	if err != nil {
		return err
	}

	err = installer.Install(p, url, name)
	if err != nil {
		return err
	}

	// TODO: write
	// dir := filepath.Dir(absPath)
	// return p.Write(dir)
	return nil
}
