package update

import (
	"github.com/RRethy/krepe/krepe/pkg/git"
	"github.com/RRethy/krepe/krepe/pkg/merger"
	"github.com/RRethy/krepe/krepe/pkg/pkg"
	"github.com/RRethy/krepe/krepe/pkg/writer"
)

type Updater struct {
	git    *git.Git
	writer writer.Writer
	merger merger.Merger
}

type Option func(*Updater)

func WithGit(git *git.Git) Option {
	return func(updater *Updater) {
		updater.git = git
	}
}

func WithWriter(w writer.Writer) Option {
	return func(updater *Updater) {
		updater.writer = w
	}
}

func WithMerger(m merger.Merger) Option {
	return func(updater *Updater) {
		updater.merger = m
	}
}

func NewUpdater(options ...Option) (*Updater, error) {
	git, err := git.NewGit()
	if err != nil {
		return nil, err
	}

	i := &Updater{
		git,
		writer.Noop{},
		merger.Noop{},
	}
	for _, option := range options {
		option(i)
	}
	return i, nil
}

func (updater *Updater) Update(p *pkg.Package, url, name string) error {
	upstreamPkgRef, err := git.NewPkgRefFromString(url)
	if err != nil {
		return err
	}

	upstreamPkgPath, err := updater.git.Clone(upstreamPkgRef)
	if err != nil {
		return err
	}

	upstreamPkg, err := pkg.NewPackageFromPathWithName(upstreamPkgPath, name)
	if err != nil {
		return err
	}

	pkgImport, err := p.GetPackageImportByName(name)
	if err != nil {
		return err
	}

	localPkgRef := pkgImport.Ref
	localPkg := pkgImport.Package

	originPkgPath, err := updater.git.Clone(localPkgRef)
	if err != nil {
		return err
	}

	originPkg, err := pkg.NewPackageFromPathWithName(originPkgPath, name)
	if err != nil {
		return err
	}

	newPkg, err := updater.merger.Merge(originPkg, localPkg, upstreamPkg)
	if err != nil {
		return err
	}

	p.UpdatePackage(newPkg, upstreamPkgRef, name)

	err = updater.writer.Write(p)
	if err != nil {
		return err
	}

	return nil
}
