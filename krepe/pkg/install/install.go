package install

import (
	"fmt"
	"path/filepath"

	"github.com/RRethy/krepe/krepe/pkg/cache"
	"github.com/RRethy/krepe/krepe/pkg/git"
	"github.com/RRethy/krepe/krepe/pkg/pkg"
)

func Install(pkgPath, url, name string) error {
	absPath, err := filepath.Abs(pkgPath)
	if err != nil {
		return fmt.Errorf("getting absolute path: %w", err)
	}
	dir := filepath.Dir(absPath)

	p, err := pkg.NewPkgFromPath(pkgPath)
	if err != nil {
		return err
	}

	installer := NewInstaller(
		git.NewGit(),
		cache.NewCache(),
	)
	err = installer.Install(p, url, name)
	if err != nil {
		return err
	}

	return p.Write(dir)
}
