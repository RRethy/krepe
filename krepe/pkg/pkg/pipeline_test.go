package pkg

import (
	"testing"

	"github.com/RRethy/krepe/krepe/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestNewPipeline(t *testing.T) {
	t.Run("valid steps", func(t *testing.T) {
		pipeline, err := NewPipeline(types.Pipeline{
			Name: "test-pipeline",
			Steps: []types.Step{
				{Function: "set-labels", Target: types.Target{Kind: "Deployment"}, Config: map[string]any{"foo": "bar"}},
				{Function: "add-annotations", Target: types.Target{APIVersion: "apps/v1"}, Config: map[string]any{"baz": "bin"}},
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, len(pipeline.Steps), 2)
		assert.Equal(t, pipeline.Steps[0].Name, "set-labels")
		assert.Equal(t, pipeline.Steps[1].Name, "add-annotations")
	})

	t.Run("invalid steps", func(t *testing.T) {
		_, err := NewPipeline(types.Pipeline{
			Name: "test-pipeline",
			Steps: []types.Step{
				{Function: "non-existent-function", Target: types.Target{Kind: "Deployment"}, Config: map[string]any{"foo": "bar"}},
			},
		})
		assert.Error(t, err)
	})
}
