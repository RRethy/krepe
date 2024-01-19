package run

import (
	"fmt"
	"path/filepath"
)

func Run(pkg, pipeline, function string) error {
	absPath, err := filepath.Abs(pkg)
	if err != nil {
		return fmt.Errorf("getting absolute path: %w", err)
	}
	dir := filepath.Dir(absPath)

	r, err := newRunnable(absPath, pipeline, function)
	if err != nil {
		return fmt.Errorf("creating runnable: %w", err)
	}

	err = r.run(dir)
	if err != nil {
		return fmt.Errorf("calling the runnable in pkg `%s`: %w", pkg, err)
	}

	return nil
}
