package install

import (
	"fmt"
	"path/filepath"

	"github.com/RRethy/krepe/krepe/pkg/git"
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

	writer, err := writer.NewDiskWriter()
	if err != nil {
		return err
	}

	git, err := git.NewGit()
	if err != nil {
		return err
	}

	installer := &Installer{
		Git:    git,
		Writer: writer,
		Dir:    absPath,
	}

	err = installer.Install(p, url, name)
	if err != nil {
		return err
	}

	return nil
}
