package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePkgRef(t *testing.T) {
	tests := []struct {
		name    string
		pkgRef  string
		want    *PkgRef
		wantErr bool
	}{
		{
			name:   "valid repoRef without path",
			pkgRef: "github.com/RRethy/krepe@v0.0.1",
			want: &PkgRef{
				Owner: "RRethy",
				Repo:  "krepe",
				Path:  nil,
				Tag:   "v0.0.1",
				Name:  "krepe",
			},
			wantErr: false,
		},
		{
			name:   "valid repoRef with path",
			pkgRef: "github.com/RRethy/krepe/pkg/repoRef@v0.0.1",
			want: &PkgRef{
				Owner: "RRethy",
				Repo:  "krepe",
				Path:  []string{"pkg", "repoRef"},
				Tag:   "v0.0.1",
				Name:  "repoRef",
			},
			wantErr: false,
		},
		{
			name:    "invalid repoRef with missing tag",
			pkgRef:  "github.com/RRethy/krepe",
			wantErr: true,
		},
		{
			name:    "invalid repoRef with missing repo",
			pkgRef:  "github.com/RRethy@v0.0.1",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := NewPkgRefFromString(test.pkgRef)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, got)
			}
		})
	}
}

func TestPkgRefString(t *testing.T) {
	tests := []struct {
		name   string
		pkgRef *PkgRef
		want   string
	}{
		{
			name: "valid repoRef without path",
			pkgRef: &PkgRef{
				Owner: "RRethy",
				Repo:  "krepe",
				Path:  nil,
				Tag:   "v0.0.1",
				Name:  "krepe",
			},
			want: "github.com/RRethy/krepe@v0.0.1",
		},
		{
			name: "valid repoRef with path",
			pkgRef: &PkgRef{
				Owner: "RRethy",
				Repo:  "krepe",
				Path:  []string{"pkg", "repoRef"},
				Tag:   "v0.0.1",
				Name:  "repoRef",
			},
			want: "github.com/RRethy/krepe/pkg/repoRef@v0.0.1",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.pkgRef.String()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestPkgRefURL(t *testing.T) {
	tests := []struct {
		name   string
		pkgRef *PkgRef
		want   string
	}{
		{
			name: "valid repoRef",
			pkgRef: &PkgRef{
				Owner: "RRethy",
				Repo:  "krepe",
				Path:  []string{"pkg", "repoRef"},
				Tag:   "v0.0.1",
				Name:  "repoRef",
			},
			want: "https://github.com/RRethy/krepe",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.pkgRef.URL()
			assert.Equal(t, test.want, got)
		})
	}
}
