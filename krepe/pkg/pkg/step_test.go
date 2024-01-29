package pkg

import (
	"reflect"
	"testing"

	"github.com/RRethy/krepe/krepe/pkg/pkg/functions"
	"github.com/RRethy/krepe/krepe/pkg/types"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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
		assert.Error(t, err)
		assert.Nil(t, step)
	})

	t.Run("invalid step with invalid target", func(t *testing.T) {
		step, err := NewStep(types.Step{Function: "set-labels", Target: types.Target{APIVersion: "apps/v1", Group: "apps"}, ConfigMap: map[string]any{"foo": "bar"}})
		assert.Error(t, err)
		assert.Nil(t, step)
	})
}

func TestStepRun(t *testing.T) {
	tests := []struct {
		name    string
		step    *types.Step
		want    *unstructured.Unstructured
		wantErr bool
	}{
		{
			name:    "valid function passes without error",
			step:    &types.Step{Function: "set-labels", Target: types.Target{Kind: "Deployment"}, ConfigMap: map[string]any{"foo": "bar2"}},
			want:    &unstructured.Unstructured{Object: map[string]any{"apiVersion": "apps/v1", "kind": "Deployment", "metadata": map[string]any{"labels": map[string]any{"foo": "bar2"}}}},
			wantErr: false,
		},
		{
			name:    "no-op if step doesn't target resource",
			step:    &types.Step{Function: "set-labels", Target: types.Target{Kind: "Service"}, ConfigMap: map[string]any{"foo": "bar2"}},
			want:    &unstructured.Unstructured{Object: map[string]any{"apiVersion": "apps/v1", "kind": "Deployment", "metadata": map[string]any{"labels": map[string]any{"foo": "bar"}}}},
			wantErr: false,
		},
		{
			name:    "function fails",
			step:    &types.Step{Function: "run-fail", Target: types.Target{Kind: "Deployment"}, ConfigMap: map[string]any{"foo": "bar"}},
			want:    &unstructured.Unstructured{Object: map[string]any{"apiVersion": "apps/v1", "kind": "Deployment", "metadata": map[string]any{"labels": map[string]any{"foo": "bar"}}}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileImport, err := NewFileImportFromBytes("deployment.yaml", []byte("apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  labels:\n    foo: bar\n"))
			assert.NoError(t, err)

			step, err := NewStep(*tt.step)
			assert.NoError(t, err)

			err = step.Run(fileImport.Resource)
			assert.Equal(t, tt.want, fileImport.Resource)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
