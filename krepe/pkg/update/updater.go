package update

import (
	"github.com/RRethy/krepe/krepe/pkg/git"
	"github.com/RRethy/krepe/krepe/pkg/pkg"
	"github.com/RRethy/krepe/krepe/pkg/pkg/imports"
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

func (updater *Updater) Update(p *pkg.Pkg, url, name string) error {
	ref, err := git.NewPkgRefFromString(url)
	if err != nil {
		return err
	}

	newPkgPath, err := updater.git.Clone(ref)
	if err != nil {
		return err
	}

	newPkg, err := pkg.NewPkgFromPath(newPkgPath)
	if err != nil {
		return err
	}

	pkgImport := imports.NewPkg(ref, name)
	err = p.UpdatePackage(newPkg, pkgImport)
	if err != nil {
		return err
	}

	return nil
}
