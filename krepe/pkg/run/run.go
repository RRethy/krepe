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

	w, err := writer.NewDiskWriter()
	if err != nil {
		return fmt.Errorf("creating writer: %w", err)
	}

	p := &Pipeline{Pkg: pkg, Name: pipeline, Writer: w, Dir: absPath}
	err = p.Run()
	if err != nil {
		return fmt.Errorf("calling the runnable in pkg `%s`: %w", pkgPath, err)
	}

	return nil
}
