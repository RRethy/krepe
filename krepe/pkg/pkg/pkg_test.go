package pkg

import (
	"path/filepath"
	"testing"

	"github.com/RRethy/krepe/krepe/pkg/git"
	"github.com/RRethy/krepe/krepe/pkg/pkg/imports"
	"github.com/RRethy/krepe/krepe/pkg/pkg/pipeline"
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.pkg.RunPipelineByName(test.pipeline)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				test.validate(test.pkg)
			}
		})
	}
}

func TestPkgRunPipeline(t *testing.T) {
	pkg, err := NewPkgFromPath(samplePkgPath)
	assert.NoError(t, err)
	goodPipeline, ok := pkg.Krepe.Pipelines.Get("mypipeline")
	assert.True(t, ok)
	badPipeline, ok := pkg.Krepe.Pipelines.Get("badpipeline")
	assert.True(t, ok)

	tests := []struct {
		name     string
		pkg      *Pkg
		validate func(*Pkg)
		pipeline pipeline.Pipeline
		wantErr  bool
	}{
		{
			name: "succeeds with valid pipeline",
			pkg:  pkg,
			validate: func(pkg *Pkg) {
				assert.Equal(t, "bar", pkg.resources[0].GetLabels()["foo"])
				assert.Equal(t, "bar", pkg.resources[1].GetLabels()["foo"])
				assert.Equal(t, "bar", pkg.resources[2].GetLabels()["foo"])
			},
			pipeline: goodPipeline,
			wantErr:  false,
		},
		{
			name:     "fails with bad pipeline",
			pkg:      pkg,
			validate: func(pkg *Pkg) {},
			pipeline: badPipeline,
			wantErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.pkg.RunPipeline(test.pipeline)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				test.validate(test.pkg)
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
			name:       "fails with duplicate pkg add",
			pkgPath:    samplePkgWithPkgInstalledPath,
			newPkgPath: samplePkgPath,
			newPkgName: "sample_pkg",
			wantErr:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pkg, err := NewPkgFromPath(test.pkgPath)
			assert.NoError(t, err)
			newPkg, err := NewPkgFromPath(test.newPkgPath)
			assert.NoError(t, err)
			repoRef, err := git.NewPkgRefFromString("github.com/RRethy/" + test.newPkgName + "@v0.0.1")
			assert.NoError(t, err)
			pkgImport := imports.NewPkg(repoRef, "")

			err = pkg.AddPackage(newPkg, pkgImport)
			if test.wantErr {
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repoRef, err := git.NewPkgRefFromString(test.other)
			assert.NoError(t, err)
			pkgImport := imports.NewPkg(repoRef, "")
			assert.Equal(t, test.want, pkg.ContainsPkg(pkgImport))
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
