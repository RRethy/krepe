package git

import (
	"bytes"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/RRethy/krepe/krepe/pkg/cache"
	"github.com/RRethy/krepe/krepe/pkg/exec"
	"github.com/stretchr/testify/assert"
)

func TestNewGit(t *testing.T) {
	cache := cache.NewCache()
	cachePath, err := cache.Path()
	assert.Nil(t, err)
	e := exec.NewExec(
		exec.WithCmd("git"),
		exec.WithDir(cachePath),
	)
	g, err := NewGit()
	assert.Nil(t, err)
	assert.NotNil(t, g)
	assert.Equal(
		t,
		e,
		g.executable,
	)
	assert.Equal(
		t,
		cachePath,
		g.dir,
	)
}

func TestGitClone(t *testing.T) {
	tmpDir := t.TempDir()
	out := &bytes.Buffer{}
	e := exec.NewExec(
		exec.WithCmd("echo"),
		exec.WithDir(tmpDir),
		exec.WithStdouterr(out),
	)

	pkgRef, err := NewPkgRefFromString("github.com/Foo/Bar/baz@v1.0.0")
	assert.Nil(t, err)

	g, err := NewGit(
		WithExec(e),
		WithDir(tmpDir),
	)
	assert.Nil(t, err)

	path, err := g.Clone(pkgRef)
	assert.Nil(t, err)
	assert.Equal(
		t,
		filepath.Join(tmpDir, "Bar", "baz"),
		path,
	)
	assert.Equal(
		t,
		fmt.Sprintf(`-C %s/Bar clone https://github.com/Foo/Bar . --depth 1
-C %s/Bar rev-parse v1.0.0
-C %s/Bar checkout v1.0.0
`, tmpDir, tmpDir, tmpDir),
		out.String(),
	)
}
