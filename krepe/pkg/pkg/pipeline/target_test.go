package pipeline

import (
	"testing"

	"github.com/RRethy/krepe/krepe/pkg/pkg/resource"
	"github.com/RRethy/krepe/krepe/pkg/yaml"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestTargetUnmarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		want    *Target
		wantErr bool
	}{
		{
			name: "succeeds with valid target",
			yaml: `
kind: Deployment
name: foo
`,
			want: &Target{
				Kind: "Deployment",
				Name: "foo",
			},
			wantErr: false,
		},
		{
			name: "success with valid apiVersion",
			yaml: `
apiVersion: apps/v1
kind: Deployment
`,
			want: &Target{
				APIVersion: "apps/v1",
				Group:      "apps",
				Version:    "v1",
				Kind:       "Deployment",
			},
			wantErr: false,
		},
		{
			name: "fails with invalid apiVersion",
			yaml: `
apiVersion: foo/bar/baz
kind: Deployment
`,
			want:    nil,
			wantErr: true,
		},
		{
			name: "fails with apiVersion and group",
			yaml: `
apiVersion: apps/v1
group: apps
kind: Deployment
`,
			want:    nil,
			wantErr: true,
		},
		{
			name: "fails with apiVersion and version",
			yaml: `
apiVersion: apps/v1
version: v1
kind: Deployment
`,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &Target{}
			err := yaml.Unmarshal([]byte(tt.yaml), got)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestTargetMarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		target  *Target
		wantYml string
		wantErr bool
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := yaml.Marshal(tt.target)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantYml, string(got))
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.target.Matches(&resource.Resource{
				Unstructured: unstructured.Unstructured{
					Object: tt.resObject,
				},
			})
			assert.Equal(t, tt.want, got)
		})
	}
}
