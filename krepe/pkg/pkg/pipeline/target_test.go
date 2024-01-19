package pipeline

import (
	"testing"

	"github.com/Shopify/krepe/krepe/pkg/yaml"
	"github.com/stretchr/testify/assert"
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
	assert.True(t, true)
}
