package install

import (
	"testing"

	"github.com/RRethy/krepe/krepe/pkg/exec"
	"github.com/RRethy/krepe/krepe/pkg/git"
	"github.com/RRethy/krepe/krepe/pkg/pkg"
	"github.com/stretchr/testify/assert"
)

const (
	samplePkgPath                 = "../../testdata/packages/sample_pkg"
	nonExistentPkgPath            = "../../testdata/packages/non_existent_pkg"
	samplePkgNameWithPkgInstalled = "../../testdata/packages/sample_pkg_with_pkg_installed"
	pacakgesDirName               = "../../testdata/packages"
)

func TestInstallerInstall(t *testing.T) {
	installedRepoRef, err := git.NewPkgRefFromString("github.com/RRethy/sample_pkg@v0.0.1")
	assert.Nil(t, err)

	tests := []struct {
		testName string
		pkgPath  string
		url      string
		name     string
		cmd      string
		wantPkg  *pkg.PackageImport
		wantErr  bool
	}{
		{
			testName: "valid install",
			pkgPath:  samplePkgPath,
			url:      "github.com/RRethy/sample_pkg@v0.0.1",
			name:     "sample_pkg",
			cmd:      "true",
			wantPkg:  pkg.NewPackageImport(installedRepoRef, "sample_pkg"),
			wantErr:  false,
		},
		{
			testName: "invalid url",
			pkgPath:  samplePkgPath,
			url:      "github.com/RRethy",
			name:     "",
			cmd:      "true",
			wantErr:  true,
		},
		{
			testName: "clone fails",
			pkgPath:  samplePkgPath,
			url:      "github.com/RRethy/sample_pkg@v0.0.1",
			name:     "",
			cmd:      "false",
			wantErr:  true,
		},
		{
			testName: "invalid pkg",
			pkgPath:  samplePkgPath,
			url:      "github.com/RRethy/non_existent_pkg@v0.0.1",
			name:     "",
			cmd:      "true",
			wantErr:  true,
		},
		{
			testName: "duplicate pkg",
			pkgPath:  samplePkgNameWithPkgInstalled,
			url:      "github.com/RRethy/sample_pkg@v0.0.1",
			name:     "",
			cmd:      "true",
			wantPkg:  pkg.NewPackageImport(installedRepoRef, ""),
			wantErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			pkg, err := pkg.NewPkgFromPath(test.pkgPath)
			assert.Nil(t, err)

			echo := exec.NewExec(exec.WithCmd(test.cmd))
			git, err := git.NewGit(
				git.WithExec(echo),
				git.WithDir(pacakgesDirName),
			)
			assert.Nil(t, err)

			installer, err := NewInstaller(WithGit(git))
			assert.Nil(t, err)

			err = installer.Install(pkg, test.url, test.name)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, pkg.ContainsPkg(test.wantPkg))
			}
		})
	}
}
