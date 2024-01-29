package install

import (
	"github.com/RRethy/krepe/krepe/pkg/git"
	"github.com/RRethy/krepe/krepe/pkg/pkg"
)

type Installer struct {
	git *git.Git
}

type Option func(*Installer)

func WithGit(git *git.Git) Option {
	return func(installer *Installer) {
		installer.git = git
	}
}

func NewInstaller(options ...Option) (*Installer, error) {
	git, err := git.NewGit()
	if err != nil {
		return nil, err
	}

	i := &Installer{
		git,
	}
	for _, option := range options {
		option(i)
	}
	return i, nil
}

func (installer *Installer) Install(p *pkg.Package, url, name string) error {
	// ref, err := git.NewPkgRefFromString(url)
	// if err != nil {
	// 	return err
	// }
	//
	// newPkgPath, err := installer.git.Clone(ref)
	// if err != nil {
	// 	return err
	// }
	//
	// newPkg, err := pkg.NewPkgFromPath(newPkgPath)
	// if err != nil {
	// 	return err
	// }
	//
	// pkgImport := pkg.NewPackageImport(ref, name)
	// err = p.AddPackage(newPkg, pkgImport)
	// if err != nil {
	// 	return err
	// }

	return nil
}
