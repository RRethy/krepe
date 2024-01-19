package pkg

import (
	"path/filepath"
	"testing"

	"github.com/RRethy/krepe/krepe/pkg/git"
	"github.com/RRethy/krepe/krepe/pkg/pkg/imports"
	"github.com/stretchr/testify/assert"
)

const (
	samplePkgPath                 = "../../testdata/packages/sample_pkg"
	samplePkgWithPkgInstalledPath = "../../testdata/packages/sample_pkg_with_pkg_installed"
)

func TestNewPkgFromPath(t *testing.T) {
	pkg, err := NewPkgFromPath(samplePkgPath)
	assert.NoError(t, err)
	assert.NotNil(t, pkg)
	assert.Equal(t, "sample_pkg", pkg.name)
	assert.Equal(t, "krepe.io/v1, Kind=Krepe", pkg.Krepe.GroupVersionKind().String())
	var fnames []string
	for _, r := range pkg.resources {
		fnames = append(fnames, r.Fname())
	}
	assert.Equal(
		t,
		[]string{"deployment.yaml", "service.yaml", "ingress.yaml"},
		fnames,
	)
	assert.Nil(t, pkg.packages)
}

func TestPkgRunPipelineByName(t *testing.T) {
	pkg, _ := NewPkgFromPath(samplePkgPath)

	tests := []struct {
		name     string
		pkg      *Pkg
		validate func(*Pkg)
		pipeline string
		wantErr  bool
	}{
		{
			name: "succeeds with valid pipeline name",
			pkg:  pkg,
			validate: func(pkg *Pkg) {
				assert.Equal(t, "bar", pkg.resources[0].GetLabels()["foo"])
				assert.Equal(t, "bar", pkg.resources[1].GetLabels()["foo"])
				assert.Equal(t, "bar", pkg.resources[2].GetLabels()["foo"])
			},
			pipeline: "mypipeline",
			wantErr:  false,
		},
		{
			name:     "fails with invalid pipeline name",
			pkg:      pkg,
			validate: func(pkg *Pkg) {},
			pipeline: "invalid",
			wantErr:  true,
		},
		{
			name:     "fails with bad pipeline",
			pkg:      pkg,
			validate: func(pkg *Pkg) {},
			pipeline: "badpipeline",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.pkg.RunPipelineByName(tt.pipeline)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				tt.validate(tt.pkg)
			}
		})
	}
}

func TestPkgAddPackage(t *testing.T) {
	tests := []struct {
		name       string
		pkgPath    string
		newPkgPath string
		newPkgName string
		wantErr    bool
	}{
		{
			name:       "succeeds with valid pkg add",
			pkgPath:    samplePkgPath,
			newPkgPath: samplePkgPath,
			newPkgName: "sample_pkg",
			wantErr:    false,
		},
		{
			name:       "fails with invalid pkg add",
			pkgPath:    samplePkgPath,
			newPkgPath: samplePkgPath,
			newPkgName: "foo",
			wantErr:    true,
		},
		{
			name:       "fails with duplicate pkg add",
			pkgPath:    samplePkgWithPkgInstalledPath,
			newPkgPath: samplePkgPath,
			newPkgName: "sample_pkg",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg, err := NewPkgFromPath(tt.pkgPath)
			assert.NoError(t, err)
			newPkg, err := NewPkgFromPath(tt.newPkgPath)
			assert.NoError(t, err)
			repoRef, err := git.NewRepoRefFromString("github.com/RRethy/" + tt.newPkgName + "@v0.0.1")
			assert.NoError(t, err)
			pkgImport := imports.NewPkg(repoRef, "")

			err = pkg.AddPackage(newPkg, pkgImport)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, pkg.ContainsPkg(pkgImport))
			}
		})
	}
}

func TestPkgContainsPkg(t *testing.T) {
	pkg, err := NewPkgFromPath(samplePkgWithPkgInstalledPath)
	assert.NoError(t, err)

	tests := []struct {
		name  string
		other string
		want  bool
	}{
		{
			name:  "contains pkg",
			other: "github.com/RRethy/sample_pkg@v0.0.1",
			want:  true,
		},
		{
			name:  "does not contain pkg",
			other: "github.com/RRethy/non_existent_pkg@v0.0.1",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoRef, err := git.NewRepoRefFromString(tt.other)
			assert.NoError(t, err)
			pkgImport := imports.NewPkg(repoRef, "")
			assert.Equal(t, tt.want, pkg.ContainsPkg(pkgImport))
		})
	}
}

func TestPkgWrite(t *testing.T) {
	tmpDir := t.TempDir()
	pkg, _ := NewPkgFromPath(samplePkgPath)
	err := pkg.Write(tmpDir)
	assert.NoError(t, err)

	got, _ := NewPkgFromPath(filepath.Join(tmpDir, pkg.name))
	assert.Equal(t, pkg.name, got.name)
	assert.Equal(t, pkg.Krepe, got.Krepe)
	var gotResources []string
	for _, r := range got.resources {
		gotResources = append(gotResources, r.Fname())
	}
	var wantResources []string
	for _, r := range pkg.resources {
		wantResources = append(wantResources, r.Fname())
	}
	assert.Equal(t, wantResources, gotResources)
	var gotPackages []string
	for _, p := range got.packages {
		gotPackages = append(gotPackages, p.name)
	}
	var wantPackages []string
	for _, p := range pkg.packages {
		wantPackages = append(wantPackages, p.name)
	}
	assert.Equal(t, wantPackages, gotPackages)
}
