package reader

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/RRethy/krepe/krepe/pkg/pkg"
	"github.com/stretchr/testify/assert"
)

func TestCacheRead(t *testing.T) {
	tests := []struct {
		name    string
		pkgPath string
		wantErr bool
		assert  func(t *testing.T, pkg *pkg.Package)
	}{
		{
			name:    "valid package cache",
			pkgPath: filepath.Join(testPackagePath, "sample_pkg_with_cache"),
			wantErr: false,
			assert: func(t *testing.T, pkg *pkg.Package) {
				assert.Equal(t, "sample_pkg", pkg.Name)
				assert.Equal(t, 3, len(pkg.FileImports))
				assert.Equal(t, 0, len(pkg.PackageImports))
				assert.Equal(t, 1, len(pkg.Pipelines))
			},
		},
		{
			name:    "valid package with subpackages cache",
			pkgPath: filepath.Join(testPackagePath, "sample_pkg_with_pkg_installed_pkg_with_cache"),
			wantErr: false,
			assert: func(t *testing.T, pkg *pkg.Package) {
				assert.Equal(t, "sample_pkg_with_pkg_installed_pkg_with_cache", pkg.Name)
				assert.Equal(t, 3, len(pkg.FileImports))
				assert.Equal(t, 1, len(pkg.PackageImports))
				assert.Equal(t, 1, len(pkg.Pipelines))
				assert.Equal(t, 3, len(pkg.PackageImports[0].Package.FileImports))
				assert.Equal(t, 0, len(pkg.PackageImports[0].Package.PackageImports))
				assert.Equal(t, 5, len(pkg.PackageImports[0].Package.Pipelines))
			},
		},
		{
			name:    "valid package with missing subpackage cache",
			pkgPath: filepath.Join(testPackagePath, "sample_pkg_with_pkg_installed_missing_cache"),
			wantErr: false,
			assert: func(t *testing.T, pkg *pkg.Package) {
				assert.Equal(t, "sample_pkg_with_pkg_installed_missing_cache", pkg.Name)
				assert.Equal(t, 3, len(pkg.FileImports))
				assert.Equal(t, 0, len(pkg.PackageImports))
				assert.Equal(t, 0, len(pkg.Pipelines))
			},
		},
		{
			name:    "no package cache",
			pkgPath: filepath.Join(testPackagePath, "sample_pkg"),
			wantErr: false,
			assert: func(t *testing.T, pkg *pkg.Package) {
				assert.Nil(t, pkg)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := Cache{}
			pkg, err := c.Read(test.pkgPath)
			if os.IsNotExist(err) {
				err = nil
			}

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				test.assert(t, pkg)
			}
		})
	}
}
