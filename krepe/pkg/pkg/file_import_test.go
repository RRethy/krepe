package pkg

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	deploymentFile = "../../testdata/packages/sample_pkg/deployment.yaml"
)

func TestNewFileImportFromPath(t *testing.T) {
	r, err := NewFileImportFromPath(deploymentFile)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "deployment.yaml", r.Filename)
	raw, err := os.ReadFile(deploymentFile)
	assert.NoError(t, err)
	assert.Equal(t, raw, r.Raw)
	assert.Equal(t, "apps/v1, Kind=Deployment", r.Resource.GroupVersionKind().String())
}

func TestNewFileImportFromBytes(t *testing.T) {
	raw, err := os.ReadFile(deploymentFile)
	assert.NoError(t, err)
	r, err := NewFileImportFromBytes("deployment.yaml", raw)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "deployment.yaml", r.Filename)
	assert.Equal(t, raw, r.Raw)
	assert.Equal(t, "apps/v1, Kind=Deployment", r.Resource.GroupVersionKind().String())
}

func TestFileImportFname(t *testing.T) {
	r, err := NewFileImportFromPath(deploymentFile)
	assert.NoError(t, err)
	assert.Equal(t, "deployment.yaml", r.Filename)
}
