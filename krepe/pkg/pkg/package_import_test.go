package pkg

import (
	"testing"

	"github.com/RRethy/krepe/krepe/pkg/git"
	"github.com/RRethy/krepe/krepe/pkg/yaml"
	"github.com/stretchr/testify/assert"
)

func TestPackageImportUnmarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		want    *PackageImport
		wantErr bool
	}{
		{
			name: "valid yaml",
			yaml: `ref: github.com/RRethy/krepe@v0.0.1
name: foo`,
			want: &PackageImport{
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
			pkg := &PackageImport{}
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

func TestPackageImportMarshalYAML(t *testing.T) {
	want := `ref: github.com/RRethy/krepe@v0.0.1
name: foo
`
	pkg := &PackageImport{
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

func TestPackageImportName(t *testing.T) {
	tests := []struct {
		name string
		pkg  *PackageImport
		want string
	}{
		{
			name: "name set",
			pkg: &PackageImport{
				Ref: &git.PkgRef{
					Owner: "RRethy",
					Repo:  "krepe",
					Path:  nil,
					Tag:   "v0.0.1",
					Name:  "krepe",
				},
				name: "foo",
			},
			want: "foo",
		},
		{
			name: "name not set",
			pkg: &PackageImport{
				Ref: &git.PkgRef{
					Owner: "RRethy",
					Repo:  "krepe",
					Path:  nil,
					Tag:   "v0.0.1",
					Name:  "krepe",
				},
				name: "",
			},
			want: "krepe",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.pkg.Name())
		})
	}
}
