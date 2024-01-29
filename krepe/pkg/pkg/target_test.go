package pkg

import (
	"testing"

	"github.com/RRethy/krepe/krepe/pkg/types"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestNewTarget(t *testing.T) {
	tests := []struct {
		name    string
		target  types.Target
		wantErr bool
		want    Target
	}{
		{
			name:    "valid target",
			target:  types.Target{Group: "core", Version: "v1", Kind: "Pod"},
			wantErr: false,
			want:    Target{Group: "core", Version: "v1", Kind: "Pod"},
		},
		{
			name:    "apiVersion and group",
			target:  types.Target{APIVersion: "apps/v1", Group: "core"},
			wantErr: true,
		},
		{
			name:    "invalid apiVersion",
			target:  types.Target{APIVersion: "a/b/c"},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewTarget(test.target)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, test.want)
			}
		})
	}
}

func TestTargetMatches(t *testing.T) {
	tests := []struct {
		name      string
		target    *Target
		resObject map[string]any
		want      bool
	}{
		{
			name: "Test Match All Fields",
			target: &Target{
				APIVersion: "v1",
				Group:      "",
				Version:    "v1",
				Kind:       "Pod",
				Name:       "test-pod",
				Namespace:  "default",
			},
			resObject: map[string]any{
				"apiVersion": "v1",
				"kind":       "Pod",
				"metadata": map[string]any{
					"name":      "test-pod",
					"namespace": "default",
				},
			},
			want: true,
		},
		{
			name: "Test Match one field",
			target: &Target{
				APIVersion: "v1",
			},
			resObject: map[string]any{
				"apiVersion": "v1",
				"kind":       "Pod",
				"metadata": map[string]any{
					"name":      "test-pod",
					"namespace": "default",
				},
			},
			want: true,
		},
		{
			name: "Test No Match APIVersion",
			target: &Target{
				APIVersion: "v2",
				Group:      "core",
				Version:    "v1",
				Kind:       "Pod",
				Name:       "test-pod",
				Namespace:  "default",
			},
			resObject: map[string]any{
				"apiVersion": "v1",
				"kind":       "Pod",
				"metadata": map[string]any{
					"name":      "test-pod",
					"namespace": "default",
				},
			},
			want: false,
		},
		{
			name: "Test No Match Group",
			target: &Target{
				APIVersion: "v1",
				Group:      "apps",
				Version:    "v1",
				Kind:       "Pod",
				Name:       "test-pod",
				Namespace:  "default",
			},
			resObject: map[string]any{
				"apiVersion": "v1",
				"kind":       "Pod",
				"metadata": map[string]any{
					"name":      "test-pod",
					"namespace": "default",
				},
			},
			want: false,
		},
		{
			name: "Test No Match Version",
			target: &Target{
				APIVersion: "v1",
				Group:      "core",
				Version:    "v2",
				Kind:       "Pod",
				Name:       "test-pod",
				Namespace:  "default",
			},
			resObject: map[string]any{
				"apiVersion": "v1",
				"kind":       "Pod",
				"metadata": map[string]any{
					"name":      "test-pod",
					"namespace": "default",
				},
			},
			want: false,
		},
		{
			name: "Test No Match Kind",
			target: &Target{
				APIVersion: "v1",
				Group:      "core",
				Version:    "v1",
				Kind:       "Service",
				Name:       "test-pod",
				Namespace:  "default",
			},
			resObject: map[string]any{
				"apiVersion": "v1",
				"kind":       "Pod",
				"metadata": map[string]any{
					"name":      "test-pod",
					"namespace": "default",
				},
			},
			want: false,
		},
		{
			name: "Test No Match Name",
			target: &Target{
				APIVersion: "v1",
				Group:      "core",
				Version:    "v1",
				Kind:       "Pod",
				Name:       "different-pod",
				Namespace:  "default",
			},
			resObject: map[string]any{
				"apiVersion": "v1",
				"kind":       "Pod",
				"metadata": map[string]any{
					"name":      "test-pod",
					"namespace": "default",
				},
			},
			want: false,
		},
		{
			name: "Test No Match Namespace",
			target: &Target{
				APIVersion: "v1",
				Group:      "core",
				Version:    "v1",
				Kind:       "Pod",
				Name:       "test-pod",
				Namespace:  "different-namespace",
			},
			resObject: map[string]any{
				"apiVersion": "v1",
				"kind":       "Pod",
				"metadata": map[string]any{
					"name":      "test-pod",
					"namespace": "default",
				},
			},
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.target.Matches(
				&unstructured.Unstructured{
					Object: test.resObject,
				},
			)
			assert.Equal(t, test.want, got)
		})
	}
}
