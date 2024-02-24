package pkg

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/RRethy/krepe/krepe/pkg/yaml"
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
			packageName:    "sample_pkg_with_pkg_installed_pkg",
			wantErr:        false,
			wantName:       "sample_pkg_with_pkg_installed_pkg",
			packageImports: []string{"sample_pkg"},
			fileImports:    []string{"deployment.yaml", "service.yaml", "ingress.yaml"},
			pipelines:      []string{"mypipeline", "badpipeline", "no-op-pipeline", "fail-on-file-import", "fail-on-package-import", "same-as-subpkg-pipeline"},
		},
		{
			name:        "bad package path",
			packagePath: "/",
			packageName: "foobar",
			wantErr:     true,
		},
		{
			name:        "empty package path",
			packagePath: "",
			packageName: "",
			wantErr:     true,
		},
		{
			name:        "package with krepe file with unknown fields must error",
			packagePath: filepath.Join(packagesPath, "sample_pkg_with_unknown_field"),
			packageName: "",
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
		{
			name:           "different name",
			packagePath:    filepath.Join(packagesPath, "sample_pkg"),
			packageName:    "foobar",
			wantErr:        false,
			wantName:       "foobar",
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
	tests := []struct {
		name         string
		pipelineName string
		wantErr      bool
		wantFound    bool
		assert       func(t *testing.T, p *Package)
	}{
		{
			name:         "valid pipeline",
			pipelineName: "mypipeline",
			wantErr:      false,
			wantFound:    true,
			assert: func(t *testing.T, p *Package) {
				assert.Equal(t, map[string]string{"foo": "bar"}, p.FileImports[0].Resource.GetLabels())
			},
		},
		{
			name:         "pipeline that exists in root and package import runs on package import first",
			pipelineName: "same-as-subpkg-pipeline",
			wantFound:    true,
			assert: func(t *testing.T, p *Package) {
				assert.Equal(t, map[string]string{"foo": "rootpkg"}, p.FileImports[2].Resource.GetLabels())
				assert.Equal(t, map[string]string{"foo": "subpkg"}, p.PackageImports[0].Package.FileImports[0].Resource.GetLabels())
				assert.Equal(t, map[string]string{"foo": "subpkg"}, p.PackageImports[0].Package.FileImports[1].Resource.GetLabels())
				assert.Equal(t, map[string]string{"foo": "rootpkg"}, p.PackageImports[0].Package.FileImports[2].Resource.GetLabels())
			},
		},
		{
			name:         "pipeline that exists only in package import runs",
			pipelineName: "only-in-subpkg-pipeline",
			wantFound:    true,
			assert: func(t *testing.T, p *Package) {
				assert.Equal(t, map[string]string{"foo": "subpkg"}, p.PackageImports[0].Package.FileImports[0].Resource.GetLabels())
				assert.Equal(t, map[string]string{"foo": "subpkg"}, p.PackageImports[0].Package.FileImports[1].Resource.GetLabels())
				assert.Equal(t, map[string]string{"foo": "subpkg"}, p.PackageImports[0].Package.FileImports[2].Resource.GetLabels())
			},
		},
		{
			name:         "non-existent pipeline",
			pipelineName: "non-existent-pipeline",
			wantErr:      false,
			wantFound:    false,
		},
		{
			name:         "fails on bad pipeline",
			pipelineName: "badpipeline",
			wantErr:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pkg, err := NewPackageFromPath(filepath.Join(packagesPath, "sample_pkg_with_pkg_installed_pkg"))
			assert.NoError(t, err)

			found, err := pkg.RunPipelineByName(test.pipelineName)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, test.wantFound, found)
				assert.NoError(t, err)
				if test.wantFound {
					test.assert(t, pkg)
				}
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
		name       string
		pkgPath    string
		newPkgPath string
		newPkgName string
		wantErr    bool
	}{
		{
			name:       "valid package",
			pkgPath:    filepath.Join(packagesPath, "sample_pkg"),
			newPkgPath: "../sample_pkg",
			newPkgName: "sample_pkg",
			wantErr:    false,
		},
		{
			name:       "valid package with name",
			pkgPath:    filepath.Join(packagesPath, "sample_pkg"),
			newPkgPath: "../sample_pkg",
			newPkgName: "foobar",
			wantErr:    false,
		},
		{
			name:       "duplicate pacakge",
			pkgPath:    filepath.Join(packagesPath, "sample_pkg_with_pkg_installed_pkg"),
			newPkgPath: "../sample_pkg",
			newPkgName: "",
			wantErr:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p, err := NewPackageFromPath(test.pkgPath)
			assert.NoError(t, err)

			newP, err := NewPackageFromPathWithName(filepath.Join(test.pkgPath, test.newPkgPath), test.newPkgName)
			assert.NoError(t, err)

			originPkgImportsLen := len(p.PackageImports)
			err = p.AddPackage(newP, test.newPkgPath)
			if test.wantErr {
				assert.Error(t, err)
				assert.Equal(t, originPkgImportsLen, len(p.PackageImports))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, originPkgImportsLen+1, len(p.PackageImports))
				assert.Equal(t, test.newPkgName, p.PackageImports[len(p.PackageImports)-1].Package.Name)
			}
		})
	}
}

func TestPackageGetPackageImportByName(t *testing.T) {
	tests := []struct {
		name    string
		pkgPath string
		pkgName string
		wantOk  bool
	}{
		{
			name:    "valid package",
			pkgPath: filepath.Join(packagesPath, "sample_pkg_with_pkg_installed_pkg"),
			pkgName: "sample_pkg",
			wantOk:  true,
		},
		{
			name:    "non-existent package",
			pkgPath: filepath.Join(packagesPath, "sample_pkg_with_pkg_installed_pkg"),
			pkgName: "non-existent-pkg",
			wantOk:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p, err := NewPackageFromPath(test.pkgPath)
			assert.NoError(t, err)

			pkgImport, ok := p.GetPackageImportByName(test.pkgName)
			if test.wantOk {
				assert.True(t, ok)
				assert.Equal(t, test.pkgName, pkgImport.Package.Name)
			} else {
				assert.False(t, ok)
			}
		})
	}
}

func TestPackageUpdatePackage(t *testing.T) {
	tests := []struct {
		name       string
		pkgPath    string
		newPkgPath string
		newPkgName string
		wantErr    bool
		wantIdx    int
		wantName   string
	}{
		{
			name:       "existent package",
			pkgPath:    filepath.Join(packagesPath, "sample_pkg_with_pkg_installed_pkg"),
			newPkgPath: "../sample_pkg",
			newPkgName: "",
			wantIdx:    0,
			wantName:   "sample_pkg",
		},
		{
			name:       "non-existent package",
			pkgPath:    filepath.Join(packagesPath, "sample_pkg_with_pkg_installed_pkg"),
			newPkgPath: "../sample_pkg",
			newPkgName: "foobar",
			wantErr:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p, err := NewPackageFromPath(test.pkgPath)
			assert.NoError(t, err)

			newP, err := NewPackageFromPathWithName(filepath.Join(test.pkgPath, test.newPkgPath), test.newPkgName)
			assert.NoError(t, err)

			p.UpdatePackage(newP, test.newPkgPath)
			if test.wantErr {
				assert.Equal(t, 1, len(p.PackageImports))
			} else {
				assert.Equal(t, test.wantName, p.PackageImports[test.wantIdx].Package.Name)
			}
		})
	}
}

func TestPackageGetTypesKrepe(t *testing.T) {
	pkg, err := NewPackageFromPath(filepath.Join(packagesPath, "sample_pkg"))
	assert.NoError(t, err)

	krepe := pkg.GetTypesKrepe()
	assert.NotNil(t, krepe)

	got, err := yaml.Marshal(krepe)
	assert.NoError(t, err)
	expected, err := os.ReadFile(filepath.Join(packagesPath, "sample_pkg", "krepe.yaml"))
	assert.NoError(t, err)
	assert.Equal(t, string(expected), string(got))
}
