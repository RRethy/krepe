package pkg

import (
	"testing"

	"github.com/Shopify/krepe/cli/pkg/pkg/imports"
	"github.com/Shopify/krepe/cli/pkg/pkg/pipeline"
	"github.com/stretchr/testify/assert"
)

const (
	krepeFile            = "../../testdata/sample_pkg/krepe.yaml"
	badKrepeFile         = "../../testdata/bad_krepe_file_pkg/krepe.yaml"
	nonExistentKrepeFile = "../../testdata/non_existent_pkg/krepe.yaml"
)

func TestNewKrepeFromPath(t *testing.T) {
	t.Run("succeeds with valid krepe file", func(t *testing.T) {
		k, err := NewKrepeFromPath(krepeFile)
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
	})

	t.Run("fails with invalid krepe file", func(t *testing.T) {
		k, err := NewKrepeFromPath(badKrepeFile)
		assert.Error(t, err)
		assert.Nil(t, k)
	})

	t.Run("fails with non-existent krepe file", func(t *testing.T) {
		k, err := NewKrepeFromPath(nonExistentKrepeFile)
		assert.Error(t, err)
		assert.Nil(t, k)
	})
}
