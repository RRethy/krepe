package exec

import (
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

	ls := NewExec(
		WithCmd("ls"),
		WithDir(tmpDir),
	)
	out, err := ls.Run()
	assert.NoError(t, err)
	assert.Equal(t, "bar.txt\nfoo.txt\n", string(out))
}
