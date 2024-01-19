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
	upstreamRef, err := git.NewPkgRefFromString(url)
	if err != nil {
		return err
	}

	upstreamPkgPath, err := updater.git.Clone(upstreamRef)
	if err != nil {
		return err
	}

	upstreamPkg, err := pkg.NewPkgFromPath(upstreamPkgPath)
	if err != nil {
		return err
	}

	upstreamPkgImport := imports.NewPkg(upstreamRef, name)

	originPkgImport := p.GetPkgImport(upstreamPkgImport.Name())
	var originPkg *pkg.Pkg
	if originPkgImport != nil {
		originPkgPath, err := updater.git.Clone(originPkgImport.Ref)
		if err != nil {
			return err
		}

		originPkg, err = pkg.NewPkgFromPath(originPkgPath)
		if err != nil {
			return err
		}
	}

	err = p.UpdatePackage(originPkg, upstreamPkg, upstreamPkgImport)
	if err != nil {
		return err
	}

	return nil
}
