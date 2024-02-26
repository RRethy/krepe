package reader

import (
	"path/filepath"
	"testing"

	"github.com/RRethy/krepe/krepe/pkg/pkg"
	"github.com/stretchr/testify/assert"
)

const (
	testPackagePath = "../../testdata/packages"
)

func TestDiskRead(t *testing.T) {
	tests := []struct {
		name    string
		pkgPath string
		wantErr bool
		assert  func(t *testing.T, pkg *pkg.Package)
	}{
		{
			name:    "valid package with no subpackages",
			pkgPath: filepath.Join(testPackagePath, "sample_pkg"),
			wantErr: false,
			assert: func(t *testing.T, pkg *pkg.Package) {
				assert.Equal(t, "sample_pkg", pkg.Name)
				assert.Equal(t, 3, len(pkg.FileImports))
				assert.Equal(t, 0, len(pkg.PackageImports))
				assert.Equal(t, 3, len(pkg.Pipelines))
			},
		},
		{
			name:    "valid package with subpackages",
			pkgPath: filepath.Join(testPackagePath, "sample_pkg_with_pkg_installed_pkg"),
			wantErr: false,
			assert: func(t *testing.T, pkg *pkg.Package) {
				assert.Equal(t, "sample_pkg_with_pkg_installed_pkg", pkg.Name)
				assert.Equal(t, 3, len(pkg.FileImports))
				assert.Equal(t, 1, len(pkg.PackageImports))
				assert.Equal(t, 6, len(pkg.Pipelines))
				assert.Equal(t, 3, len(pkg.PackageImports[0].Package.FileImports))
				assert.Equal(t, 3, len(pkg.PackageImports[0].Package.FileImports))
				assert.Equal(t, 0, len(pkg.PackageImports[0].Package.PackageImports))
				assert.Equal(t, 5, len(pkg.PackageImports[0].Package.Pipelines))
			},
		},
		{
			name:    "non-existent package path",
			pkgPath: filepath.Join(testPackagePath, "zzzzz"),
			wantErr: true,
		},
		{
			name:    "non-directory package path",
			pkgPath: filepath.Join(testPackagePath, "sample_pkg", "krepe.yaml"),
			wantErr: true,
		},
		{
			name:    "bad imported file",
			pkgPath: filepath.Join(testPackagePath, "bad_imported_file_pkg"),
			wantErr: true,
		},
		{
			name:    "bad imported package",
			pkgPath: filepath.Join(testPackagePath, "bad_imported_package_pkg"),
			wantErr: true,
		},
		{
			name:    "bad imported pipeline",
			pkgPath: filepath.Join(testPackagePath, "bad_pipeline_pkg"),
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d := Disk{}
			pkg, err := d.Read(test.pkgPath)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				test.assert(t, pkg)
			}
		})
	}
}
