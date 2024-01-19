package git

import (
	"os"

	"github.com/RRethy/krepe/krepe/pkg/reporef"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type Git struct {
}

func NewGit() *Git {
	return &Git{}
}

func (g *Git) Clone(ref reporef.RepoRef, path string) error {
	cloned := true
	var err error
	if _, err = os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		cloned = false
	}

	var repo *git.Repository
	if !cloned {
		repo, err = git.PlainClone(path, false, &git.CloneOptions{
			URL:      ref.URL,
			Progress: os.Stdout,
		})
		if err != nil {
			return err
		}
	}

	wt, err := repo.Worktree()
	if err != nil {
		return err
	}

	err = wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName("refs/tags/" + ref.Tag),
	})
	if err != nil {
		return err
	}

	return nil
}
