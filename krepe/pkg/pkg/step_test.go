package pkg

import (
	"reflect"
	"testing"

	"github.com/RRethy/krepe/krepe/pkg/pkg/functions"
	"github.com/RRethy/krepe/krepe/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestNewStep(t *testing.T) {
	t.Run("valid step", func(t *testing.T) {
		step, err := NewStep(types.Step{Function: "set-labels", Target: types.Target{Kind: "Deployment"}, ConfigMap: map[string]any{"foo": "bar"}})
		assert.NoError(t, err)
		assert.Equal(t, "set-labels", step.Name)
		assert.Equal(t, reflect.TypeOf(&functions.SetLabels{}), reflect.TypeOf(step.Fn))
		assert.Equal(t, Target{Kind: "Deployment"}, step.Target)
		assert.Equal(t, map[string]any{"foo": "bar"}, step.ConfigMap)
	})

	t.Run("invalid step with non-existent function", func(t *testing.T) {
		step, err := NewStep(types.Step{Function: "zzz", Target: types.Target{Kind: "Deployment"}, ConfigMap: map[string]any{"foo": "bar"}})
		assert.NoError(t, err)
		assert.Nil(t, step)
	})

	t.Run("invalid step with invalid target", func(t *testing.T) {
		step, err := NewStep(types.Step{Function: "set-labels", Target: types.Target{APIVersion: "apps/v1", Group: "apps"}, ConfigMap: map[string]any{"foo": "bar"}})
		assert.NoError(t, err)
		assert.Nil(t, step)
	})
}

func TestStepRun(t *testing.T) {
	resource, err := NewResourceFromBytes("deployment.yaml", []byte("apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  labels:\n    foo: bar\n"))
	assert.NoError(t, err)

	t.Run("valid function passes without error", func(t *testing.T) {
		step, err := NewStep(types.Step{Function: "set-labels", Target: types.Target{Kind: "Deployment"}, ConfigMap: map[string]any{"foo": "bar2"}})
		assert.NoError(t, err)

		err = step.Run(resource)
		assert.NoError(t, err)
		assert.Equal(t, "bar2", resource.GetLabels()["foo"])
	})

	t.Run("no-op if step doesn't target resource", func(t *testing.T) {
		step, err := NewStep(types.Step{Function: "set-labels", Target: types.Target{Kind: "Service"}, ConfigMap: map[string]any{"foo": "bar2"}})
		assert.NoError(t, err)

		err = step.Run(resource)
		assert.NoError(t, err)
		assert.Equal(t, "bar", resource.GetLabels()["foo"])
	})

	t.Run("function fails", func(t *testing.T) {
		step, err := NewStep(types.Step{Function: "run-fail", Target: types.Target{Kind: "Deployment"}, ConfigMap: map[string]any{"foo": "bar"}})
		assert.NoError(t, err)

		err = step.Run(resource)
		assert.Error(t, err)
		assert.Equal(t, []byte("apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  labels:\n    foo: bar\n"), resource.Raw)
	})
}
