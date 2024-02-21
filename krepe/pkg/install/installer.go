package install

import (
	"github.com/RRethy/krepe/krepe/pkg/git"
	"github.com/RRethy/krepe/krepe/pkg/pkg"
	"github.com/RRethy/krepe/krepe/pkg/writer"
)

type Installer struct {
	Git    *git.Git
	Writer writer.Writer
	Dir    string
}

func (installer *Installer) Install(p *pkg.Package, url, name string) error {
	ref, err := git.NewPkgRefFromString(url)
	if err != nil {
		return err
	}

	newPkgPath, err := installer.Git.Clone(ref)
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

	err = installer.Writer.Write(p, installer.Dir)
	if err != nil {
		return err
	}

	return nil
}
