package install

import (
	"fmt"
	"path/filepath"

	"github.com/RRethy/krepe/krepe/pkg/pkg"
	"github.com/RRethy/krepe/krepe/pkg/writer"
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

	writer, err := writer.NewPackageWriter(filepath.Dir(absPath))
	if err != nil {
		return err
	}

	installer, err := NewInstaller(
		WithWriter(writer),
	)
	if err != nil {
		return err
	}

	err = installer.Install(p, url, name)
	if err != nil {
		return err
	}

	return nil
}
