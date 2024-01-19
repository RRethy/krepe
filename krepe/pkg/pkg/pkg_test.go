package pkg

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	samplePkgPath = "../../testdata/packages/sample_pkg"
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

func TestInstallPackage(t *testing.T) {
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
