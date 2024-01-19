package imports

import (
	"testing"

	"github.com/RRethy/krepe/krepe/pkg/git"
	"github.com/RRethy/krepe/krepe/pkg/yaml"
	"github.com/stretchr/testify/assert"
)

func TestPkgUnmarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		want    *Pkg
		wantErr bool
	}{
		{
			name: "valid yaml",
			yaml: `ref: github.com/RRethy/krepe@v0.0.1
name: foo`,
			want: &Pkg{
				ref: &git.RepoRef{
					URL:  "github.com/RRethy/krepe",
					Tag:  "v0.0.1",
					Name: "krepe",
				},
				name: "foo",
			},
			wantErr: false,
		},
		{
			name:    "invalid yaml",
			yaml:    `ref`,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid pkg import",
			yaml:    `ref: github.com/RRethy/krepe`,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg := &Pkg{}
			err := yaml.Unmarshal([]byte(tt.yaml), pkg)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, pkg)
			}
		})
	}
}

func TestPkgMarshalYAML(t *testing.T) {
	want := `ref: github.com/RRethy/krepe@v0.0.1
name: foo
`
	pkg := &Pkg{
		ref: &git.RepoRef{
			URL: "github.com/RRethy/krepe",
			Tag: "v0.0.1",
		},
		name: "foo",
	}
	got, err := yaml.Marshal(&pkg)
	assert.NoError(t, err)
	assert.Equal(t, want, string(got))
}
