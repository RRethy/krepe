package install

import (
	"fmt"
	"path/filepath"

	"github.com/RRethy/krepe/krepe/pkg/cache"
	"github.com/RRethy/krepe/krepe/pkg/git"
	"github.com/RRethy/krepe/krepe/pkg/pkg"
	"github.com/RRethy/krepe/krepe/pkg/pkg/imports"
)

type Installer struct {
	gitClient   git.Client
	cacheClient *cache.Client
}

func NewInstaller(git git.Client, cache *cache.Client) *Installer {
	return &Installer{
		gitClient:   git,
		cacheClient: cache,
	}
}

func (i *Installer) Install(pkgPath, url, name string) error {
	absPath, err := filepath.Abs(pkgPath)
	if err != nil {
		return fmt.Errorf("getting absolute path: %w", err)
	}
	dir := filepath.Dir(absPath)

	p, err := pkg.NewPkgFromPath(pkgPath)
	if err != nil {
		return err
	}

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

	err = p.Write(dir)
	if err != nil {
		return err
	}

	return nil
}
