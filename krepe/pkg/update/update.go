package update

import (
	"fmt"
	"path/filepath"

	"github.com/RRethy/krepe/krepe/pkg/merger"
	"github.com/RRethy/krepe/krepe/pkg/pkg"
	"github.com/RRethy/krepe/krepe/pkg/writer"
)

func Update(pkgPath, url, name string) error {
	absPath, err := filepath.Abs(pkgPath)
	if err != nil {
		return fmt.Errorf("getting absolute path: %w", err)
	}

	p, err := pkg.NewPackageFromPath(pkgPath)
	if err != nil {
		return err
	}

	merger := merger.NewMerger()

	writer, err := writer.NewDiskWriter()
	if err != nil {
		return err
	}

	updater, err := NewUpdater(WithMerger(merger), WithWriter(writer), WithDir(absPath))
	if err != nil {
		return err
	}

	err = updater.Update(p, url, name)
	if err != nil {
		return err
	}

	return nil
}
