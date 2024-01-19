package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	samplePkgPath = "../../testdata/sample_pkg"
)

func TestNewPkgFromPath(t *testing.T) {
	pkg, err := NewPkgFromPath(samplePkgPath)
	assert.NoError(t, err)
	assert.NotNil(t, pkg)
	assert.Equal(t, "sample_pkg", pkg.name)
	assert.Equal(t, "krepe.io/v1, Kind=Krepe", pkg.krepe.GroupVersionKind().String())
	var fnames []string
	for _, r := range pkg.resources {
		fnames = append(fnames, r.Fname())
	}
	assert.Equal(
		t,
		[]string{"deployment.yaml", "service.yaml", "ingress.yaml"},
		fnames,
	)
	assert.Nil(t, pkg.packages)
}
