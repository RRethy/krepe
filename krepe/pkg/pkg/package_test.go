package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	samplePkgPath                 = "../../testdata/packages/sample_pkg"
	samplePkgWithPkgInstalledPath = "../../testdata/packages/sample_pkg_with_pkg_installed"
)

func TestNewPackageFromPath(t *testing.T) {
	p, err := NewPackageFromPath(samplePkgPath)
	assert.NoError(t, err)
	assert.Equal(t, "sample_pkg", p.Name)
}

func TestNewPackageFromPathWithName(t *testing.T) {
}

func TestPackageRunPipelineByName(t *testing.T) {
}

func TestPackageRunPipeline(t *testing.T) {
}
