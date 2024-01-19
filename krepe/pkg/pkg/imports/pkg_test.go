package imports

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPkg(t *testing.T) {
	tests := []struct {
		name    string
		tag     string
		pkgName string
		want    *Pkg
		wantErr bool
	}{
		{
			name:    "valid tag with empty name",
			tag:     "github.com/Shopify/krepe@v0.0.1",
			pkgName: "",
			want: &Pkg{
				url:     "github.com/Shopify/krepe",
				version: "v0.0.1",
				name:    "krepe",
			},
			wantErr: false,
		},
		{
			name:    "valid tag with valid name",
			tag:     "github.com/Shopify/krepe@v0.0.1",
			pkgName: "notkrepe",
			want: &Pkg{
				url:     "github.com/Shopify/krepe",
				version: "v0.0.1",
				name:    "notkrepe",
			},
			wantErr: false,
		},
		{
			name:    "invalid tag with no @",
			tag:     "github.com/Shopify/krepe",
			pkgName: "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid tag with no url",
			tag:     "@v0.0.1",
			pkgName: "",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg, err := NewPkg(tt.tag, tt.pkgName)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, pkg)
			}
		})
	}
}

func TestPkgUnmarshalYAML(t *testing.T) {
}

func TestPkgMarshalYAML(t *testing.T) {
}
