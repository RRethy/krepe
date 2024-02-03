package pkg

import (
	"path/filepath"
	"testing"

	"github.com/RRethy/krepe/krepe/pkg/git"
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
		wantName       string
		packageImports []string
		fileImports    []string
		pipelines      []string
	}{
		{
			name:           "valid package",
			packagePath:    filepath.Join(packagesPath, "sample_pkg_with_pkg_installed_pkg"),
			packageName:    "sample_pkg_with_pkg_installed",
			wantErr:        false,
			wantName:       "sample_pkg_with_pkg_installed",
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
			name:        "bad name",
			packagePath: filepath.Join(packagesPath, "sample_pkg"),
			packageName: "/",
			wantErr:     true,
		},
		{
			name:           "empty name",
			packagePath:    filepath.Join(packagesPath, "sample_pkg"),
			packageName:    "",
			wantErr:        false,
			wantName:       "sample_pkg",
			packageImports: nil,
			fileImports:    []string{"deployment.yaml", "service.yaml", "ingress.yaml"},
			pipelines:      []string{"mypipeline", "badpipeline", "no-op-pipeline"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p, err := NewPackageFromPathWithName(test.packagePath, test.packageName)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.wantName, p.Name)

				assert.Equal(t, len(test.packageImports), len(p.PackageImports))
				assert.Equal(t, len(test.fileImports), len(p.FileImports))
				assert.Equal(t, len(test.pipelines), len(p.Pipelines))

				for i, pkgImport := range p.PackageImports {
					assert.Equal(t, test.packageImports[i], pkgImport.Package.Name)
				}
				for i, fileImport := range p.FileImports {
					assert.Equal(t, test.fileImports[i], fileImport.Name)
				}
				for i, pipeline := range p.Pipelines {
					assert.Equal(t, test.pipelines[i], pipeline.Name)
				}
			}
		})
	}
}

func TestPackageRunPipelineByName(t *testing.T) {
	pkg, err := NewPackageFromPath(filepath.Join(packagesPath, "sample_pkg_with_pkg_installed_pkg"))
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
				// assert.Equal(t, map[string]string{"foo": "bar"}, p.FileImports[0].Resource.GetLabels())
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
				// test.assert(t, pkg)
			}
		})
	}
}

func TestPackageRunPipeline(t *testing.T) {
	pkg, err := NewPackageFromPath(filepath.Join(packagesPath, "sample_pkg_with_pkg_installed_pkg"))
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

func TestPackageAddPackage(t *testing.T) {
	tests := []struct {
		name         string
		pkgPath      string
		newPkgPath   string
		newPkgRefStr string
		newPkgName   string
		wantErr      bool
	}{
		{
			name:         "valid package",
			pkgPath:      filepath.Join(packagesPath, "sample_pkg"),
			newPkgPath:   filepath.Join(packagesPath, "sample_pkg"),
			newPkgRefStr: "github.com/RRethy/sample_pkg@v0.0.1",
			newPkgName:   "",
			wantErr:      false,
		},
		{
			name:         "valid package with name",
			pkgPath:      filepath.Join(packagesPath, "sample_pkg"),
			newPkgPath:   filepath.Join(packagesPath, "sample_pkg"),
			newPkgRefStr: "github.com/RRethy/sample_pkg@v0.0.1",
			newPkgName:   "foobar",
			wantErr:      false,
		},
		{
			name:         "duplicate pacakge",
			pkgPath:      filepath.Join(packagesPath, "sample_pkg_with_pkg_installed_pkg"),
			newPkgPath:   filepath.Join(packagesPath, "sample_pkg"),
			newPkgRefStr: "github.com/RRethy/sample_pkg@v0.0.1",
			newPkgName:   "",
			wantErr:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p, err := NewPackageFromPath(test.pkgPath)
			assert.NoError(t, err)

			newP, err := NewPackageFromPath(test.newPkgPath)
			assert.NoError(t, err)

			ref, err := git.NewPkgRefFromString(test.newPkgRefStr)
			assert.NoError(t, err)

			originPkgImportsLen := len(p.PackageImports)
			err = p.AddPackage(newP, ref, test.newPkgName)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, originPkgImportsLen+1, len(p.PackageImports))
				expectedName := ref.Name
				if test.newPkgName != "" {
					expectedName = test.newPkgName
				}
				assert.Equal(t, expectedName, p.PackageImports[len(p.PackageImports)-1].Package.Name)
			}
		})
	}
}

func TestPackageGetPackageImportByName(t *testing.T) {
	tests := []struct {
		name    string
		pkgPath string
		pkgName string
		wantErr bool
	}{
		{
			name:    "valid package",
			pkgPath: filepath.Join(packagesPath, "sample_pkg_with_pkg_installed_pkg"),
			pkgName: "sample_pkg",
			wantErr: false,
		},
		{
			name:    "non-existent package",
			pkgPath: filepath.Join(packagesPath, "sample_pkg_with_pkg_installed_pkg"),
			pkgName: "non-existent-pkg",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p, err := NewPackageFromPath(test.pkgPath)
			assert.NoError(t, err)

			pkgImport, err := p.GetPackageImportByName(test.pkgName)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.pkgName, pkgImport.Package.Name)
			}
		})
	}
}

func TestPackageUpdatePackage(t *testing.T) {
	tests := []struct {
		name         string
		pkgPath      string
		newPkgPath   string
		newPkgRefStr string
		newPkgName   string
		wantIdx      int
		wantName     string
	}{
		{
			name:         "existent package",
			pkgPath:      filepath.Join(packagesPath, "sample_pkg_with_pkg_installed_pkg"),
			newPkgPath:   filepath.Join(packagesPath, "sample_pkg"),
			newPkgRefStr: "github.com/RRethy/sample_pkg@v0.0.1",
			newPkgName:   "",
			wantIdx:      0,
			wantName:     "sample_pkg",
		},
		{
			name:         "non-existent package",
			pkgPath:      filepath.Join(packagesPath, "sample_pkg_with_pkg_installed_pkg"),
			newPkgPath:   filepath.Join(packagesPath, "sample_pkg"),
			newPkgRefStr: "github.com/RRethy/sample_pkg@v0.0.1",
			newPkgName:   "foobar",
			wantIdx:      1,
			wantName:     "foobar",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p, err := NewPackageFromPath(test.pkgPath)
			assert.NoError(t, err)

			newP, err := NewPackageFromPath(test.newPkgPath)
			assert.NoError(t, err)

			ref, err := git.NewPkgRefFromString(test.newPkgRefStr)
			assert.NoError(t, err)

			p.UpdatePackage(newP, ref, test.newPkgName)
			assert.Equal(t, test.wantName, p.PackageImports[test.wantIdx].Package.Name)
		})
	}
}
