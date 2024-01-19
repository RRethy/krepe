package update

import (
	"github.com/RRethy/krepe/krepe/pkg/git"
	"github.com/RRethy/krepe/krepe/pkg/merger"
	"github.com/RRethy/krepe/krepe/pkg/pkg"
)

type Updater struct {
	git    *git.Git
	merger merger.Merger[*pkg.Pkg]
}

type Option func(*Updater)

func WithGit(git *git.Git) Option {
	return func(updater *Updater) {
		updater.git = git
	}
}

func WithMerger(merger merger.Merger[*pkg.Pkg]) Option {
	return func(updater *Updater) {
		updater.merger = merger
	}
}

func NewUpdater(options ...Option) (*Updater, error) {
	git, err := git.NewGit()
	if err != nil {
		return nil, err
	}

	i := &Updater{
		git,
		&merger.PkgMerger{},
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

	upstreamPkgImport := pkg.NewPackageImport(upstreamRef, name)
	pkgName := upstreamPkgImport.Name()

	originPkgImport, ok := p.GetPkgImport(pkgName)
	var originPkg *pkg.Pkg
	if ok {
		originPkgPath, err := updater.git.Clone(originPkgImport.Ref)
		if err != nil {
			return err
		}

		originPkg, err = pkg.NewPkgFromPath(originPkgPath)
		if err != nil {
			return err
		}
	}

	localPkg, _ := p.GetPkg(pkgName)
	newPkg := updater.merger.ThreeWayMerge(originPkg, localPkg, upstreamPkg)
	err = p.UpdatePackage(newPkg, upstreamPkgImport)
	if err != nil {
		return err
	}

	return nil
}
