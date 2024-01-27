package update

import (
	"github.com/RRethy/krepe/krepe/pkg/git"
	_ "github.com/RRethy/krepe/krepe/pkg/merger"
	"github.com/RRethy/krepe/krepe/pkg/pkg"
)

type Updater struct {
	git *git.Git
}

type Option func(*Updater)

func WithGit(git *git.Git) Option {
	return func(updater *Updater) {
		updater.git = git
	}
}

func NewUpdater(options ...Option) (*Updater, error) {
	git, err := git.NewGit()
	if err != nil {
		return nil, err
	}

	i := &Updater{
		git,
	}
	for _, option := range options {
		option(i)
	}
	return i, nil
}

func (updater *Updater) Update(p *pkg.Package, url, name string) error {
	return nil
	// upstreamRef, err := git.NewPkgRefFromString(url)
	// if err != nil {
	// 	return err
	// }
	//
	// upstreamPkgPath, err := updater.git.Clone(upstreamRef)
	// if err != nil {
	// 	return err
	// }
	//
	// upstreamPkg, err := pkg.NewPkgFromPath(upstreamPkgPath)
	// if err != nil {
	// 	return err
	// }
	//
	// upstreamPkgImport := pkg.NewPackageImport(upstreamRef, name)
	// pkgName := upstreamPkgImport.Name()
	//
	// originPkgImport, ok := p.GetPkgImport(pkgName)
	// var originPkg *pkg.Pkg
	// if ok {
	// 	originPkgPath, err := updater.git.Clone(originPkgImport.Ref)
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	originPkg, err = pkg.NewPkgFromPath(originPkgPath)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	//
	// localPkg, _ := p.GetPkg(pkgName)
	// newPkg, err := merger.ThreeWayMerge(originPkg, localPkg, upstreamPkg)
	// if err != nil {
	// 	return err
	// }
	//
	// // TODO: ensure name is set correctly
	// err = p.UpdatePackage(newPkg, upstreamPkgImport)
	// if err != nil {
	// 	return err
	// }
	//
	// return nil
}
