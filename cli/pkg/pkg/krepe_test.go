package pkg

import (
	"testing"

	"github.com/Shopify/krepe/cli/pkg/pkg/imports"
	"github.com/Shopify/krepe/cli/pkg/pkg/pipeline"
	"github.com/stretchr/testify/assert"
)

func TestNewKrepeFromPath(t *testing.T) {
	k, err := NewKrepeFromPath("../../testdata/sample_pkg/krepe.yaml")
	assert.NoError(t, err)
	assert.NotNil(t, k)
	assert.Equal(t, &imports.Imports{
		Files: []string{
			"deployment.yaml",
			"service.yaml",
			"ingress.yaml",
		},
		Packages: nil,
	}, k.Imports)
	assert.Equal(t, map[string]*pipeline.Pipeline(nil), k.Pipelines)
}
