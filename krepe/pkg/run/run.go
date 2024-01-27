package run

import (
	_ "fmt"
	_ "path/filepath"

	_ "github.com/RRethy/krepe/krepe/pkg/pkg"
)

func Run(pkgPath, pipeline string) error {
	return nil
	// absPath, err := filepath.Abs(pkgPath)
	// if err != nil {
	// 	return fmt.Errorf("getting absolute path: %w", err)
	// }
	// dir := filepath.Dir(absPath)
	//
	// pkg, err := pkg.NewPkgFromPath(pkgPath)
	// if err != nil {
	// 	return err
	// }
	//
	// r, err := newPipeline(pkg, pipeline), nil
	// if err != nil {
	// 	return fmt.Errorf("creating runnable: %w", err)
	// }
	//
	// err = r.run(dir)
	// if err != nil {
	// 	return fmt.Errorf("calling the runnable in pkg `%s`: %w", pkgPath, err)
	// }
	//
	// return nil
}
