package install

import (
	"path/filepath"
	"testing"

	"github.com/RRethy/krepe/krepe/pkg/writer"
	"github.com/stretchr/testify/assert"
)

const (
	packagesDirName = "../../testdata/packages"
)

func TestInstallerInstall(t *testing.T) {
	tests := []struct {
		name                   string
		pkgPath                string
		newPkgPath             string
		newPkgName             string
		wantErr                bool
		wantWriteErr           bool
		wantImportRelativePath string
		wantImportName         string
	}{
		{
			name:                   "valid install",
			pkgPath:                filepath.Join(packagesDirName, "sample_pkg"),
			newPkgPath:             "../sample_pkg_with_pkg_installed_pkg",
			newPkgName:             "foobar",
			wantErr:                false,
			wantImportRelativePath: "../sample_pkg_with_pkg_installed_pkg",
			wantImportName:         "foobar",
		},
		{
			name:                   "valid install with no name",
			pkgPath:                filepath.Join(packagesDirName, "sample_pkg"),
			newPkgPath:             "../sample_pkg_with_pkg_installed_pkg",
			newPkgName:             "",
			wantErr:                false,
			wantImportRelativePath: "../sample_pkg_with_pkg_installed_pkg",
			wantImportName:         "sample_pkg_with_pkg_installed_pkg",
		},
		{
			name:       "invalid pkg",
			pkgPath:    filepath.Join(packagesDirName, "sample_pkg"),
			newPkgPath: "../invalid_pkg",
			newPkgName: "",
			wantErr:    true,
		},
		{
			name:       "duplicate pkg",
			pkgPath:    filepath.Join(packagesDirName, "sample_pkg_with_pkg_installed_pkg"),
			newPkgPath: "../sample_pkg",
			newPkgName: "sample_pkg",
			wantErr:    true,
		},
		{
			name:         "write fails",
			pkgPath:      filepath.Join(packagesDirName, "sample_pkg"),
			newPkgPath:   "../sample_pkg_with_pkg_installed_pkg",
			newPkgName:   "foobar",
			wantErr:      true,
			wantWriteErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			writer := &writer.Mock{Success: !test.wantWriteErr}
			installer := &Installer{Writer: writer}
			err := installer.Install(test.pkgPath, test.newPkgPath, test.newPkgName)
			if test.wantErr {
				assert.Error(t, err)
				assert.Equal(t, 0, writer.Cnt)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, 1, writer.Cnt)
				pkg := writer.Pkg
				importedPackage, ok := pkg.GetPackageImportByName(test.wantImportName)
				assert.True(t, ok)
				assert.Equal(t, test.wantImportRelativePath, importedPackage.RelativePath)
				assert.Equal(t, test.wantImportName, importedPackage.Package.Name)
			}
		})
	}
}
