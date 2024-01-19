package git

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/RRethy/krepe/krepe/pkg/cache"
	"github.com/RRethy/krepe/krepe/pkg/exec"
)

type Git struct {
	executable *exec.Exec
	dir        string
}

type Option func(*Git)

func WithDir(dir string) Option {
	return func(g *Git) {
		g.dir = dir
	}
}

func WithExec(executable *exec.Exec) Option {
	return func(g *Git) {
		g.executable = executable
	}
}

func NewGit(options ...Option) (*Git, error) {
	cache := cache.NewCache()
	cachePath, err := cache.Path()
	if err != nil {
		return nil, err
	}

	g := &Git{
		executable: exec.NewExec(
			exec.WithCmd("git"),
			exec.WithDir(cachePath),
		),
		dir: cachePath,
	}
	for _, option := range options {
		option(g)
	}
	return g, nil
}

func (g *Git) Clone(ref *PkgRef) (string, error) {
	cloned := true
	clonePath := filepath.Join(g.dir, ref.Repo)
	fmt.Println(clonePath)
	var err error
	if _, err = os.Stat(clonePath); err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}

		err = os.MkdirAll(clonePath, 0755)
		if err != nil {
			return "", err
		}
		cloned = false
	}

	if !cloned {
		_, err := g.executable.Run("-C", clonePath, "clone", ref.URL(), ".", "--depth", "1")
		if err != nil {
			return "", err
		}
	}

	_, err = g.executable.Run("-C", clonePath, "rev-parse", ref.Tag)
	if err != nil {
		_, err := g.executable.Run("-C", clonePath, "fetch", "origin", "tag", ref.Tag)
		if err != nil {
			return "", err
		}
	}

	_, err = g.executable.Run("-C", clonePath, "checkout", ref.Tag)
	if err != nil {
		return "", err
	}

	return filepath.Join(append([]string{clonePath}, ref.Path...)...), nil
}
