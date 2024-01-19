package exec

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecRun(t *testing.T) {
	tmpDir := t.TempDir()
	err := exec.Command("touch", filepath.Join(tmpDir, "foo.txt")).Run()
	assert.NoError(t, err)
	err = exec.Command("touch", filepath.Join(tmpDir, "bar.txt")).Run()
	assert.NoError(t, err)
	stdouterr := &bytes.Buffer{}

	ls := NewExec(
		WithCmd("ls"),
		WithDir(tmpDir),
		WithStdouterr(stdouterr),
	)
	out, err := ls.Run()
	assert.NoError(t, err)
	assert.Equal(t, "bar.txt\nfoo.txt\n", string(out))
	assert.Equal(t, "bar.txt\nfoo.txt\n", stdouterr.String())
}
