package resource

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	deploymentFile = "../../../testdata/sample_pkg/deployment.yaml"
)

func TestNewResourceFromPath(t *testing.T) {
	r, err := NewResourceFromPath(deploymentFile)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "deployment.yaml", r.Fname())
	raw, err := os.ReadFile(deploymentFile)
	assert.NoError(t, err)
	assert.Equal(t, raw, r.raw)
	assert.Equal(t, "apps/v1, Kind=Deployment", r.GroupVersionKind().String())
}
