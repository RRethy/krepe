package pkg

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	deploymentFile = "../../testdata/packages/sample_pkg/deployment.yaml"
)

func TestNewResourceFromPath(t *testing.T) {
	r, err := NewResourceFromPath(deploymentFile)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "deployment.yaml", r.Filename)
	raw, err := os.ReadFile(deploymentFile)
	assert.NoError(t, err)
	assert.Equal(t, raw, r.Raw)
	assert.Equal(t, "apps/v1, Kind=Deployment", r.GroupVersionKind().String())
}

func TestNewResourceFromBytes(t *testing.T) {
	raw, err := os.ReadFile(deploymentFile)
	assert.NoError(t, err)
	r, err := NewResourceFromBytes("deployment.yaml", raw)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "deployment.yaml", r.Filename)
	assert.Equal(t, raw, r.Raw)
	assert.Equal(t, "apps/v1, Kind=Deployment", r.GroupVersionKind().String())
}

func TestResourceFname(t *testing.T) {
	r, err := NewResourceFromPath(deploymentFile)
	assert.NoError(t, err)
	assert.Equal(t, "deployment.yaml", r.Filename)
}
