package pkg

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	packagesPath = "../../testdata/packages"
)

func TestNewPackageFromPath(t *testing.T) {
	p, err := NewPackageFromPath(filepath.Join(packagesPath, "sample_pkg"))
	assert.NoError(t, err)
	assert.Equal(t, "sample_pkg", p.Name)
}

func TestNewPackageFromPathWithName(t *testing.T) {
	tests := []struct {
		name           string
		packagePath    string
		packageName    string
		wantErr        bool
		packageImports []string
		fileImports    []string
		pipelines      []string
	}{
		{
			name:           "valid package",
			packagePath:    filepath.Join(packagesPath, "sample_pkg_with_pkg_installed_pkg"),
			packageName:    "sample_pkg_with_pkg_installed",
			wantErr:        false,
			packageImports: []string{"sample_pkg"},
			fileImports:    []string{"deployment.yaml", "service.yaml", "ingress.yaml"},
			pipelines:      []string{"mypipeline", "badpipeline", "no-op-pipeline", "fail-on-file-import", "fail-on-package-import"},
		},
		{
			name:        "bad package path",
			packagePath: "/",
			packageName: "foobar",
			wantErr:     true,
		},
		{
			name:        "non-existent package path",
			packagePath: filepath.Join(packagesPath, "non_existent_pkg"),
			packageName: "foobar",
			wantErr:     true,
		},
		{
			name:        "non-directory package path",
			packagePath: filepath.Join(packagesPath, "sample_pkg", "krepe.yaml"),
			packageName: "foobar",
			wantErr:     true,
		},
		{
			name:        "no krepe.yaml in package",
			packagePath: filepath.Join(packagesPath, "no_krepe_file"),
			packageName: "foobar",
			wantErr:     true,
		},
		{
			name:        "bad krepe.yaml in package",
			packagePath: filepath.Join(packagesPath, "bad_krepe_file_pkg"),
			packageName: "foobar",
			wantErr:     true,
		},
		{
			name:        "bad imported file",
			packagePath: filepath.Join(packagesPath, "bad_imported_file_pkg"),
			packageName: "foobar",
			wantErr:     true,
		},
		{
			name:        "bad pipeline",
			packagePath: filepath.Join(packagesPath, "bad_pipeline_pkg"),
			packageName: "foobar",
			wantErr:     true,
		},
		{
			name:        "empty name",
			packagePath: filepath.Join(packagesPath, "sample_pkg"),
			packageName: "",
			wantErr:     true,
		},
		{
			name:        "bad name",
			packagePath: filepath.Join(packagesPath, "sample_pkg"),
			packageName: "/",
			wantErr:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p, err := NewPackageFromPathWithName(test.packagePath, test.packageName)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.packageName, p.Name)

				assert.Equal(t, len(test.packageImports), len(p.PackageImports))
				assert.Equal(t, len(test.fileImports), len(p.FileImports))
				assert.Equal(t, len(test.pipelines), len(p.Pipelines))

				for i, pkgImport := range p.PackageImports {
					assert.Equal(t, test.packageImports[i], pkgImport.Name)
				}
				for i, fileImport := range p.FileImports {
					assert.Equal(t, test.fileImports[i], fileImport.Filename)
				}
				for i, pipeline := range p.Pipelines {
					assert.Equal(t, test.pipelines[i], pipeline.Name)
				}
			}
		})
	}
}

func TestPackageRunPipelineByName(t *testing.T) {
	pkg, err := NewPackageFromPath(filepath.Join(packagesPath, "sample_pkg_with_pkg_installed"))
	assert.NoError(t, err)

	tests := []struct {
		name         string
		pipelineName string
		wantErr      bool
		assert       func(t *testing.T, p *Package)
	}{
		{
			name:         "valid pipeline",
			pipelineName: "mypipeline",
			wantErr:      false,
			assert: func(t *testing.T, p *Package) {
				assert.Equal(t, map[string]string{"foo": "bar"}, p.FileImports[0].Resource.GetLabels())
			},
		},
		{
			name:         "non-existent pipeline",
			pipelineName: "non-existent-pipeline",
			wantErr:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := pkg.RunPipelineByName(test.pipelineName)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				test.assert(t, pkg)
			}
		})
	}
}

func TestPackageRunPipeline(t *testing.T) {
	pkg, err := NewPackageFromPath(filepath.Join(packagesPath, "sample_pkg_with_pkg_installed"))
	assert.NoError(t, err)

	tests := []struct {
		name     string
		pipeline Pipeline
		wantErr  bool
		assert   func(t *testing.T, p *Package)
	}{
		{
			name:     "valid pipeline",
			pipeline: pkg.Pipelines[0],
			wantErr:  false,
			assert: func(t *testing.T, p *Package) {
				assert.Equal(t, map[string]string{"foo": "bar"}, p.FileImports[0].Resource.GetLabels())
			},
		},
		{
			name:     "fails on file import",
			pipeline: pkg.Pipelines[3],
			wantErr:  true,
		},
		{
			name:     "fails on package import",
			pipeline: pkg.Pipelines[4],
			wantErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := pkg.RunPipeline(test.pipeline)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				test.assert(t, pkg)
			}
		})
	}
}
