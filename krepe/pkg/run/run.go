package run

import (
	"fmt"
	"path/filepath"

	"github.com/RRethy/krepe/krepe/pkg/pkg"
	"github.com/RRethy/krepe/krepe/pkg/writer"
)

func Run(pkgPath, pipeline string) error {
	absPath, err := filepath.Abs(pkgPath)
	if err != nil {
		return fmt.Errorf("getting absolute path: %w", err)
	}

	pkg, err := pkg.NewPackageFromPath(pkgPath)
	if err != nil {
		return err
	}

	w, err := writer.NewPackageWriter(absPath)
	if err != nil {
		return fmt.Errorf("creating writer: %w", err)
	}

	r := newPipeline(pkg, pipeline, withWriter(w))

	err = r.run()
	if err != nil {
		return fmt.Errorf("calling the runnable in pkg `%s`: %w", pkgPath, err)
	}

	return nil
}
