package install

import (
	"github.com/RRethy/krepe/krepe/pkg/cache"
	"github.com/RRethy/krepe/krepe/pkg/git"
)

func Install(pkgPath, url, name string) error {
	installer := NewInstaller(
		git.NewGit(),
		cache.NewClient(),
	)
	return installer.Install(pkgPath, url, name)
}
