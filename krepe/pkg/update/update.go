package update

import (
	"fmt"
	"path/filepath"

	"github.com/RRethy/krepe/krepe/pkg/pkg"
)

func Update(pkgPath, url, name string) error {
	absPath, err := filepath.Abs(pkgPath)
	if err != nil {
		return fmt.Errorf("getting absolute path: %w", err)
	}

	p, err := pkg.NewPkgFromPath(pkgPath)
	if err != nil {
		return err
	}

	updater, err := NewUpdater()
	if err != nil {
		return err
	}

	err = updater.Update(p, url, name)
	if err != nil {
		return err
	}

	dir := filepath.Dir(absPath)
	return p.Write(dir)
}
