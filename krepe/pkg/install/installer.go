package install

import (
	"path/filepath"

	"github.com/RRethy/krepe/krepe/pkg/cache"
	"github.com/RRethy/krepe/krepe/pkg/git"
	"github.com/RRethy/krepe/krepe/pkg/pkg"
	"github.com/RRethy/krepe/krepe/pkg/pkg/imports"
)

type Installer struct {
	gitClient   git.Client
	cacheClient cache.Cache
}

func NewInstaller(git git.Client, cache cache.Cache) *Installer {
	return &Installer{
		gitClient:   git,
		cacheClient: cache,
	}
}

func (i *Installer) Install(p *pkg.Pkg, url, name string) error {
	ref, err := git.NewRepoRefFromString(url)
	if err != nil {
		return err
	}

	pkgImport := imports.NewPkg(ref, name)
	installPath := filepath.Join(i.cacheClient.Path(), pkgImport.Name())

	err = i.gitClient.CloneInto(ref, installPath)
	if err != nil {
		return err
	}

	newPkg, err := pkg.NewPkgFromPath(installPath)
	if err != nil {
		return err
	}

	err = p.AddPackage(newPkg, pkgImport)
	if err != nil {
		return err
	}

	return nil
}
