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
				Ref: &git.PkgRef{
					Owner: "RRethy",
					Repo:  "krepe",
					Path:  nil,
					Tag:   "v0.0.1",
					Name:  "krepe",
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pkg := &Pkg{}
			err := yaml.Unmarshal([]byte(test.yaml), pkg)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, pkg)
			}
		})
	}
}

func TestPkgMarshalYAML(t *testing.T) {
	want := `ref: github.com/RRethy/krepe@v0.0.1
name: foo
`
	pkg := &Pkg{
		Ref: &git.PkgRef{
			Owner: "RRethy",
			Repo:  "krepe",
			Tag:   "v0.0.1",
		},
		name: "foo",
	}
	got, err := yaml.Marshal(&pkg)
	assert.NoError(t, err)
	assert.Equal(t, want, string(got))
}
