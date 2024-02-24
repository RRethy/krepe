package install

import (
	"path/filepath"
	"testing"

	"github.com/RRethy/krepe/krepe/pkg/pkg"
	"github.com/RRethy/krepe/krepe/pkg/writer"
	"github.com/stretchr/testify/assert"
)

const (
	packagesDirName = "../../testdata/packages"
)

func TestInstallerInstall(t *testing.T) {
	tests := []struct {
		name         string
		pkgPath      string
		url          string
		newPkgName   string
		cmd          string
		wantErr      bool
		wantWriteErr bool
		wantPkgRef   *git.PkgRef
		wantPkgName  string
		wantPackage  func(t *testing.T, pkg *pkg.Package) // called with the imported package
	}{
		{
			name:        "valid install",
			pkgPath:     filepath.Join(packagesDirName, "sample_pkg"),
			url:         "github.com/RRethy/sample_pkg@v0.0.1",
			newPkgName:  "foobar",
			cmd:         "true",
			wantErr:     false,
			wantPkgRef:  installedRepoRef,
			wantPkgName: "foobar",
			wantPackage: func(t *testing.T, pkg *pkg.Package) {
				assert.Equal(t, "foobar", pkg.Name)
				assert.Equal(t, 0, len(pkg.PackageImports))
				assert.Equal(t, 3, len(pkg.FileImports))
			},
		},
		{
			name:        "valid install with no name",
			pkgPath:     filepath.Join(packagesDirName, "sample_pkg"),
			url:         "github.com/RRethy/sample_pkg@v0.0.1",
			cmd:         "true",
			wantErr:     false,
			wantPkgRef:  installedRepoRef,
			wantPkgName: "sample_pkg",
			wantPackage: func(t *testing.T, pkg *pkg.Package) {
				assert.Equal(t, "sample_pkg", pkg.Name)
				assert.Equal(t, 0, len(pkg.PackageImports))
				assert.Equal(t, 3, len(pkg.FileImports))
			},
		},
		{
			name:    "invalid url",
			pkgPath: filepath.Join(packagesDirName, "sample_pkg"),
			url:     "github.com/RRethy/sample_pkg",
			cmd:     "true",
			wantErr: true,
		},
		{
			name:    "clone fails",
			pkgPath: filepath.Join(packagesDirName, "sample_pkg"),
			url:     "github.com/RRethy/sample_pkg@v0.0.1",
			cmd:     "false",
			wantErr: true,
		},
		{
			name:    "invalid pkg",
			pkgPath: filepath.Join(packagesDirName, "sample_pkg"),
			url:     "github.com/RRethy/invalid_pkg@v0.0.1",
			cmd:     "true",
			wantErr: true,
		},
		{
			name:    "duplicate pkg",
			pkgPath: filepath.Join(packagesDirName, "sample_pkg_with_pkg_installed_pkg"),
			url:     "github.com/RRethy/sample_pkg@v0.0.1",
			cmd:     "true",
			wantErr: true,
		},
		{
			name:         "write fails",
			pkgPath:      filepath.Join(packagesDirName, "sample_pkg"),
			url:          "github.com/RRethy/sample_pkg@v0.0.1",
			newPkgName:   "foobar",
			cmd:          "true",
			wantErr:      true,
			wantWriteErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pkg, err := pkg.NewPackageFromPathWithName(test.pkgPath, test.newPkgName)
			assert.Nil(t, err)

			writer := &writer.Mock{
				Success: !test.wantWriteErr,
			}

			installer := &Installer{Writer: writer}

			err = installer.Install(pkg, test.url, test.newPkgName)
			if test.wantErr {
				assert.Error(t, err)
				assert.Equal(t, 0, writer.Cnt)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, 1, writer.Cnt)
				importedPackage, err := pkg.GetPackageImportByName(test.wantPkgName)
				assert.Nil(t, err)
				assert.Equal(t, test.wantPkgRef, importedPackage.Ref)
				assert.Equal(t, test.wantPkgName, importedPackage.Package.Name)
				test.wantPackage(t, importedPackage.Package)
			}
		})
	}
}
