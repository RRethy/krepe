package install

import (
	"errors"
	"testing"

	"github.com/RRethy/krepe/krepe/pkg/git"
	"github.com/RRethy/krepe/krepe/pkg/pkg"
	"github.com/RRethy/krepe/krepe/pkg/pkg/imports"
	"github.com/stretchr/testify/assert"
)

const (
	// TODO: use testdata/packages
	samplePkgPath                 = "../../testdata/sample_pkg"
	nonExistentPkgPath            = "../../testdata/non_existent_pkg"
	samplePkgNameWithPkgInstalled = "../../testdata/sample_pkg_with_pkg_installed"
	pacakgesDirName               = "../../testdata"
)

type fakeCache struct {
	path string
}

func (f *fakeCache) Path() string {
	return f.path
}

type fakeGitClient struct {
	cloneFails bool
}

func (f *fakeGitClient) CloneInto(ref *git.RepoRef, dir string) error {
	if f.cloneFails {
		return errors.New("clone failed")
	}
	return nil
}

func newGitClient(cloneFails bool) git.Client {
	return &fakeGitClient{
		cloneFails: cloneFails,
	}
}

func TestInstallerInstall(t *testing.T) {
	installedRepoRef, err := git.NewRepoRefFromString("github.com/RRethy/sample_pkg@v0.0.1")
	assert.Nil(t, err)

	tests := []struct {
		testName   string
		pkgPath    string
		url        string
		name       string
		cloneFails bool
		wantPkg    *imports.Pkg
		wantErr    bool
	}{
		{
			testName:   "valid install",
			pkgPath:    samplePkgPath,
			url:        "github.com/RRethy/sample_pkg@v0.0.1",
			name:       "sample_pkg",
			cloneFails: false,
			wantPkg:    imports.NewPkg(installedRepoRef, "sample_pkg"),
			wantErr:    false,
		},
		{
			testName:   "invalid url",
			pkgPath:    samplePkgPath,
			url:        "github.com/RRethy",
			name:       "",
			cloneFails: false,
			wantErr:    true,
		},
		{
			testName:   "clone fails",
			pkgPath:    samplePkgPath,
			url:        "github.com/RRethy/sample_pkg@v0.0.1",
			name:       "",
			cloneFails: true,
			wantErr:    true,
		},
		{
			testName:   "invalid pkg",
			pkgPath:    samplePkgPath,
			url:        "github.com/RRethy/non_existent_pkg@v0.0.1",
			name:       "",
			cloneFails: false,
			wantErr:    true,
		},
		{
			testName:   "duplicate pkg",
			pkgPath:    samplePkgNameWithPkgInstalled,
			url:        "github.com/RRethy/sample_pkg@v0.0.1",
			name:       "",
			cloneFails: false,
			wantPkg:    imports.NewPkg(installedRepoRef, ""),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			pkg, err := pkg.NewPkgFromPath(tt.pkgPath)
			assert.Nil(t, err)

			installer := NewInstaller(
				newGitClient(tt.cloneFails),
				&fakeCache{path: pacakgesDirName},
			)
			err = installer.Install(pkg, tt.url, tt.name)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, pkg.ContainsPkg(tt.wantPkg))
			}
		})
	}
}
