package pkg

import (
	"testing"

	"github.com/RRethy/krepe/krepe/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestNewPipeline(t *testing.T) {
	t.Run("valid steps", func(t *testing.T) {
		pipeline, err := NewPipeline([]types.Step{
			types.Step{Function: "set-labels", Target: types.Target{Kind: "Deployment"}, ConfigMap: map[string]any{"foo": "bar"}},
			types.Step{Function: "add-annotations", Target: types.Target{APIVersion: "apps/v1"}, ConfigMap: map[string]any{"baz": "bin"}},
		})
		assert.NoError(t, err)
		assert.Equal(t, len(pipeline.Steps), 2)
		assert.Equal(t, pipeline.Steps[0].Name, "set-labels")
		assert.Equal(t, pipeline.Steps[1].Name, "add-annotations")
	})

	t.Run("invalid steps", func(t *testing.T) {
		pipeline, err := NewPipeline([]types.Step{
			types.Step{Function: "non-existent-function", Target: types.Target{Kind: "Deployment"}, ConfigMap: map[string]any{"foo": "bar"}},
		})
		assert.Error(t, err)
		assert.Nil(t, pipeline)
	})
}
