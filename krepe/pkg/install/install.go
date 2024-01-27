package install

import (
	_ "fmt"
	_ "path/filepath"

	_ "github.com/RRethy/krepe/krepe/pkg/pkg"
)

func Install(pkgPath, url, name string) error {
	return nil
	// absPath, err := filepath.Abs(pkgPath)
	// if err != nil {
	// 	return fmt.Errorf("getting absolute path: %w", err)
	// }
	//
	// p, err := pkg.NewPkgFromPath(pkgPath)
	// if err != nil {
	// 	return err
	// }
	//
	// installer, err := NewInstaller()
	// if err != nil {
	// 	return err
	// }
	//
	// err = installer.Install(p, url, name)
	// if err != nil {
	// 	return err
	// }
	//
	// dir := filepath.Dir(absPath)
	// return p.Write(dir)
}
