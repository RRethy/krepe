package git

import (
	"os"
)

type Client interface {
	CloneInto(ref *RepoRef, dir string) error
}

type Git struct{}

func NewGit() Client {
	return &Git{}
}

func (g *Git) CloneInto(ref *RepoRef, dir string) error {
	cloned := true
	var err error
	if _, err = os.Stat(dir); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		cloned = false
	}

	exec := NewExec(dir)

	if !cloned {
		_, err := exec.Run("clone", ref.URL, "--depth", "1")
		if err != nil {
			return err
		}
	}

	_, err = exec.Run("rev-parse", ref.Tag)
	if err != nil {
		_, err = exec.Run("fetch", "origin", "tag", ref.Tag)
		if err != nil {
			return err
		}
	}

	_, err = exec.Run("checkout", ref.Tag)
	if err != nil {
		return err
	}

	return nil
}
