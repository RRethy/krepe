package update

import (
	_ "path/filepath"

	"github.com/RRethy/krepe/krepe/pkg/merger"
	_ "github.com/RRethy/krepe/krepe/pkg/pkg"
	"github.com/RRethy/krepe/krepe/pkg/writer"
)

type Updater struct {
	Writer writer.Writer
	Merger merger.Merger
}

// TODO: we should never use filepath.Abs unless we know we have a full path
func (updater *Updater) Update(rootPkgPath, packageName string) error {
	// p, err := pkg.NewPackageFromPath(rootPkgPath)
	// if err != nil {
	// 	return err
	// }
	//
	// pkgImport, err := p.GetPackageImportByName(packageName)
	// if err != nil {
	// 	return err
	// }
	//
	// subPkgPath := filepath.Join(rootPkgPath, pkgImport.Package.Name)
	// subPkg, err := pkg.NewPackageFromPath(subPkgPath)
	// if err != nil {
	// 	return err
	// }
	//
	// upstreamPkgPath := filepath.Join(rootPkgPath, pkgImport.RelativePath)
	// upstreamPkg, err := pkg.NewPackageFromPath(upstreamPkgPath)
	// if err != nil {
	// 	return err
	// }

	// krepe update -C rootPkg subPkg1
	// krepe update -C rootPkg/subPkg1 subPkg1-1
	// krepe update -C rootPkg/subPkg1/subPkg1-1 subPkg1-1-1

	// rootPkg
	// .krepe/rootPkg
	// 	- subPkg1
	//    .krepe/subPkg1
	//    - subPkg1-1
	//      - subPkg1-1-1
	//    - subPkg1-2
	// 	- subPkg2
	//    .krepe/subPkg2
	//    - subPkg2-1
	//    - subPkg2-2
	// 	- subPkg3
	//    .krepe/subPkg3
	//
	// a package has a .krepe directory with a krepe.yaml and imported files
	// if subpackage doesn't exist and a parent directory is .krepe/, then we are an origin cache and we should look in ../<subpackage>/.krepe/
	// needs a change in the pkg package to handle this

	// p.UpdatePackage(newPkg, upstreamPkgRef, name)
	//
	// err = updater.Writer.Write(p, updater.dir)
	// if err != nil {
	// 	return err
	// }

	return nil
}
