package run

import (
	"path/filepath"
	"testing"

	"github.com/Shopify/krepe/cli/pkg/pkg"
	"github.com/stretchr/testify/assert"
)

func TestPipelineRun(t *testing.T) {
	tmpDir := t.TempDir()
	pkg, err := pkg.NewPkgFromPath(samplePkgPath)
	assert.NoError(t, err)
	pipeline := newPipeline(pkg, "no-op-pipeline")
	err = pipeline.run(tmpDir)
	assert.NoError(t, err)
	assert.DirExists(t, filepath.Join(tmpDir, "sample_pkg"))
}
