package install

import (
	"github.com/RRethy/krepe/krepe/pkg/git"
	"github.com/RRethy/krepe/krepe/pkg/pkg"
	"github.com/RRethy/krepe/krepe/pkg/writer"
)

type Installer struct {
	git    *git.Git
	writer writer.Writer
	dir    string
}

type Option func(*Installer)

func WithGit(git *git.Git) Option {
	return func(installer *Installer) {
		installer.git = git
	}
}

func WithWriter(w writer.Writer) Option {
	return func(installer *Installer) {
		installer.writer = w
	}
}

func WithDir(dir string) Option {
	return func(installer *Installer) {
		installer.dir = dir
	}
}

func NewInstaller(options ...Option) (*Installer, error) {
	git, err := git.NewGit()
	if err != nil {
		return nil, err
	}

	i := &Installer{
		git,
		writer.Noop{},
		"",
	}
	for _, option := range options {
		option(i)
	}
	return i, nil
}

func (installer *Installer) Install(p *pkg.Package, url, name string) error {
	ref, err := git.NewPkgRefFromString(url)
	if err != nil {
		return err
	}

	newPkgPath, err := installer.git.Clone(ref)
	if err != nil {
		return err
	}

	newPkg, err := pkg.NewPackageFromPathWithName(newPkgPath, name)
	if err != nil {
		return err
	}

	err = p.AddPackage(newPkg, ref, name)
	if err != nil {
		return err
	}

	err = installer.writer.Write(p, installer.dir)
	if err != nil {
		return err
	}

	return nil
}
