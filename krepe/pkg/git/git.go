package git

import (
	"os"

	"github.com/RRethy/krepe/krepe/pkg/reporef"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func Clone(ref *reporef.RepoRef, dir string) error {
	cloned := true
	var err error
	if _, err = os.Stat(dir); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		cloned = false
	}

	var repo *git.Repository
	if !cloned {
		repo, err = git.PlainClone(dir, false, &git.CloneOptions{
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
