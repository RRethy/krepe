package resource

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	deploymentFile = "../../../testdata/packages/sample_pkg/deployment.yaml"
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

func TestNewResourceFromBytes(t *testing.T) {
	raw, err := os.ReadFile(deploymentFile)
	assert.NoError(t, err)
	r, err := NewResourceFromBytes("deployment.yaml", raw)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "deployment.yaml", r.Fname())
	assert.Equal(t, raw, r.raw)
	assert.Equal(t, "apps/v1, Kind=Deployment", r.GroupVersionKind().String())
}

func TestResourceFname(t *testing.T) {
	r, err := NewResourceFromPath(deploymentFile)
	assert.NoError(t, err)
	assert.Equal(t, "deployment.yaml", r.Fname())
}

func TestResourceWrite(t *testing.T) {
	tmpDir := t.TempDir()
	r := &Resource{
		fname: "test.yaml",
		Unstructured: unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "Pod",
				"metadata": map[string]interface{}{
					"name":      "test-pod",
					"namespace": "default",
				},
			},
		},
	}

	wantYaml := "apiVersion: v1\nkind: Pod\nmetadata:\n  name: test-pod\n  namespace: default\n"

	err := r.Write(tmpDir)
	assert.NoError(t, err)

	got, err := os.ReadFile(filepath.Join(tmpDir, r.fname))
	assert.NoError(t, err)
	assert.Equal(t, wantYaml, string(got))

	r.fname = "testdir"
	err = os.Mkdir(filepath.Join(tmpDir, r.fname), 0755)
	assert.NoError(t, err)
	err = r.Write(tmpDir)
	assert.Error(t, err)
}
